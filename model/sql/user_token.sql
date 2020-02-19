-- Create user related to token table
CREATE TABLE `tbl_user_token` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '递增主键',
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户令牌',
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
