-- Create user related to file table
CREATE TABLE `tbl_user_file` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '递增主键',
  `user_name` varchar(64) NOT NULL COMMENT '用户昵称',
  `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件哈希',
  
  `status` int(10) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_file` (`user_name`, `file_sha1`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
