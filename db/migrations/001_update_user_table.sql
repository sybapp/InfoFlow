-- Migration: Update user table for improved password security and indexing
-- Date: 2024

-- 1. Add indexes for frequently queried fields
ALTER TABLE `user` ADD KEY `ix_username` (`username`) IF NOT EXISTS;
ALTER TABLE `user` ADD KEY `ix_phone` (`phone`) IF NOT EXISTS;

-- 2. Increase password column size to accommodate bcrypt hashes (60 characters)
ALTER TABLE `user` MODIFY COLUMN `password` varchar(256) NOT NULL DEFAULT '' COMMENT '密码';

-- Note: Existing MD5 passwords will still work due to backward compatibility in VerifyPassword function
-- New passwords will be hashed using bcrypt
-- Consider running a migration script to rehash all existing passwords when users login
