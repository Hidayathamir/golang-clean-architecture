create unique index users_username_active_idx 
on users (username) 
where (deleted_at is null);