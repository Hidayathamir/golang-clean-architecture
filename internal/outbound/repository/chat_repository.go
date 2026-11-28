package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockRepositoryChat.go -pkg=mock . ChatRepository
type ChatRepository interface {
	Create(ctx context.Context, db *gorm.DB, chat *entity.Chat) error
	FindThread(ctx context.Context, db *gorm.DB, chatList *entity.ChatList, userID, peerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error
	ListConversations(ctx context.Context, db *gorm.DB, conversations *entity.ConversationList, userID int64) error
	MarkThreadRead(ctx context.Context, db *gorm.DB, userID, peerID int64, readAt time.Time) (int64, error)
	CountUnread(ctx context.Context, db *gorm.DB, userID int64) (int64, error)
}

var _ ChatRepository = &ChatRepositoryImpl{}

type ChatRepositoryImpl struct {
	Cfg *config.Config
}

func NewChatRepository(cfg *config.Config) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{
		Cfg: cfg,
	}
}

func (r *ChatRepositoryImpl) Create(ctx context.Context, db *gorm.DB, chat *entity.Chat) error {
	err := db.WithContext(ctx).Create(chat).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*ChatRepositoryImpl).Create")
	}
	return nil
}

func (r *ChatRepositoryImpl) FindThread(ctx context.Context, db *gorm.DB, chatList *entity.ChatList, userID, peerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	wherePair, argPair := column.SenderID.Eq(userID)
	whereRecipient, argRecipient := column.RecipientID.Eq(peerID)
	whereReverse, argReverse := column.SenderID.Eq(peerID)
	whereSender, argSender := column.RecipientID.Eq(userID)
	query := db.WithContext(ctx).
		Where(db.Where(wherePair, argPair).Where(whereRecipient, argRecipient).Or(db.Where(whereReverse, argReverse).Where(whereSender, argSender))).
		Order(column.CreatedAt.Desc()).
		Order(column.ID.Desc()).
		Limit(limit)
	if cursorID > 0 {
		whereKeyset, argCreatedAt, argID := column.KeysetBefore(cursorCreatedAt, cursorID)
		query = query.Where(whereKeyset, argCreatedAt, argID)
	}
	err := query.Find(chatList).Error
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*ChatRepositoryImpl).FindThread")
	}
	return nil
}

// ListConversations returns one row per peer the user has exchanged chats with,
// carrying the latest chat. The subquery picks the newest chat id per peer.
func (r *ChatRepositoryImpl) ListConversations(ctx context.Context, db *gorm.DB, conversations *entity.ConversationList, userID int64) error {
	const q = `
select
    c.id,
    c.sender_id,
    c.recipient_id,
    c.content,
    c.read_at,
    c.created_at,
    c.updated_at
from chats c
join (
    select
        case when sender_id = ? then recipient_id else sender_id end as peer_id,
        max(id) as last_id
    from chats
    where deleted_at is null and (sender_id = ? or recipient_id = ?)
    group by peer_id
) t on t.last_id = c.id
where c.deleted_at is null
order by c.created_at desc`
	rows, err := db.WithContext(ctx).Raw(q, userID, userID, userID).Rows()
	if err != nil {
		return errkit.AddFuncName(err, "repository.(*ChatRepositoryImpl).ListConversations")
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		conversation := entity.Conversation{}
		err = rows.Scan(
			&conversation.LastChat.ID,
			&conversation.LastChat.SenderID,
			&conversation.LastChat.RecipientID,
			&conversation.LastChat.Content,
			&conversation.LastChat.ReadAt,
			&conversation.LastChat.CreatedAt,
			&conversation.LastChat.UpdatedAt,
		)
		if err != nil {
			return errkit.AddFuncName(err, "repository.(*ChatRepositoryImpl).ListConversations")
		}
		if conversation.LastChat.SenderID == userID {
			conversation.PeerID = conversation.LastChat.RecipientID
		} else {
			conversation.PeerID = conversation.LastChat.SenderID
		}
		conversation.LastChatAt = conversation.LastChat.CreatedAt
		*conversations = append(*conversations, conversation)
	}
	if rows.Err() != nil {
		return errkit.AddFuncName(rows.Err(), "repository.(*ChatRepositoryImpl).ListConversations")
	}
	return nil
}

func (r *ChatRepositoryImpl) MarkThreadRead(ctx context.Context, db *gorm.DB, userID, peerID int64, readAt time.Time) (int64, error) {
	whereRecipient, argRecipient := column.RecipientID.Eq(userID)
	whereSender, argSender := column.SenderID.Eq(peerID)
	whereUnread, argUnread := column.ReadAt.Eq(nil)
	result := db.WithContext(ctx).
		Model(&entity.Chat{}).
		Where(whereRecipient, argRecipient).
		Where(whereSender, argSender).
		Where(whereUnread, argUnread).
		Update(column.ReadAt.Str(), sql.NullTime{Time: readAt, Valid: true})
	if result.Error != nil {
		return 0, errkit.AddFuncName(result.Error, "repository.(*ChatRepositoryImpl).MarkThreadRead")
	}
	return result.RowsAffected, nil
}

func (r *ChatRepositoryImpl) CountUnread(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
	var total int64
	whereRecipient, argRecipient := column.RecipientID.Eq(userID)
	whereUnread, argUnread := column.ReadAt.Eq(nil)
	err := db.WithContext(ctx).
		Model(&entity.Chat{}).
		Where(whereRecipient, argRecipient).
		Where(whereUnread, argUnread).
		Count(&total).Error
	if err != nil {
		return 0, errkit.AddFuncName(err, "repository.(*ChatRepositoryImpl).CountUnread")
	}
	return total, nil
}
