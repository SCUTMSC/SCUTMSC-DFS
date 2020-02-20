-- Create file table
CREATE TABLE `tbl_file` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  
  `file_sha1` char(64) NOT NULL DEFAULT '' COMMENT '文件哈希',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名字',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_path` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件位置',
  `enable_times` int(10) NOT NULL DEFAULT 9999999 COMMENT '下载次数',
  `enable_days` int(10) NOT NULL DEFAULT 30 COMMENT '下载时间',
  `create_at` datetime DEFAULT NOW() COMMENT '创建日期',
  `update_at` datetime DEFAULT NOW() on update CURRENT_TIMESTAMP() COMMENT '更新日期',
  
  `status` int(10) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
  `ext1` int(10) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
