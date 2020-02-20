-- Create user table
CREATE TABLE `tbl_user` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '递增主键',
  
  `nick_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户密码',

  `real_name` varchar(64) DEFAULT '' COMMENT '用户真名',
  `school_number` varchar(128) DEFAULT '' COMMENT '用户学号',
  `gender` tinyint(1) DEFAULT 1 COMMENT '用户性别',
  `birthday` datetime DEFAULT '0000-00-00 00:00:00'COMMENT '用户生日',
  `campus` varchar(64) DEFAULT '' COMMENT '用户校区',
  `school` varchar(64) DEFAULT '' COMMENT '用户学院',
  `major` varchar(64) DEFAULT '' COMMENT '用户专业',
  `grade` varchar(64) DEFAULT '' COMMENT '用户年级',
  `class` varchar(64) DEFAULT '' COMMENT '用户班级',
  `dormitory` varchar(64) DEFAULT '' COMMENT '用户宿舍',
  `department` varchar(64) DEFAULT '' COMMENT '用户部门', 

  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(128) DEFAULT '' COMMENT '手机',
  `wechat` varchar(128) DEFAULT '' COMMENT '微信',
  `qq` varchar(128) DEFAULT '' COMMENT 'QQ',
  `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机是否已验证',
  `wechat_validated` tinyint(1) DEFAULT 0 COMMENT '微信是否已验证',
  `qq_validated` tinyint(1) DEFAULT 0 COMMENT 'QQ是否已验证',

  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `signin_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '登录时间',

  `profile` text COMMENT '用户属性',
  `status` int(10) NOT NULL DEFAULT '0' COMMENT '账户状态',

  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_nickname` (`nick_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
