/*
 Navicat Premium Data Transfer

 Source Server         : mysql-local
 Source Server Type    : MySQL
 Source Server Version : 50715
 Source Host           : 127.0.0.1
 Source Database       : blog

 Target Server Type    : MySQL
 Target Server Version : 50715
 File Encoding         : utf-8

 Date: 12/26/2016 00:10:01 AM
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `admin`
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin` (
  `id` int(11) NOT NULL,
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- ----------------------------
--  Records of `admin`
-- ----------------------------
BEGIN;
INSERT INTO `admin` VALUES ('1', 'baa', 'c3284d0f94606de1fd2af172aba15bf3', 'baa', '/assets/img/user.png', '2016-11-26 07:54:15', '2016-12-24 23:52:25');
COMMIT;

-- ----------------------------
--  Table structure for `content`
-- ----------------------------
DROP TABLE IF EXISTS `content`;
CREATE TABLE `content` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '摘要',
  `content` text NOT NULL COMMENT '内容',
  `create_user_id` int(10) unsigned NOT NULL COMMENT '创建者ID',
  `deleted` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否已删除，0否 1是',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='内容表';

-- ----------------------------
--  Records of `content`
-- ----------------------------
BEGIN;
INSERT INTO `content` VALUES ('1', 'Hello, Baa.', 'BaaisanexpressGowebframeworkwithrouting,middleware', '<p>Baa is an express Go web framework with routing, middleware, dependency injection, http context.</p><p><b>Getting Started</b><br></p><p>installing:</p><pre><code class=\"lang-python\">go get -u gopkg.in/baa.v1<br></code></pre><p>example:</p><pre><code>package main\r\n\r\nimport (\r\n    \"gopkg.in/baa.v1\"\r\n)\r\n\r\nfunc main() {\r\n    app := baa.New()\r\n    app.Get(\"/\", func(c *baa.Context) {\r\n        c.String(200, \"Hello World!\")\r\n    })\r\n    app.Run(\":8080\")\r\n}<br></code></pre><p><b>Features</b></p><ul><li>route support static, param, group<br></li><li>route support handler chain<br></li><li>route support static file serve<br></li><li>middleware supoort handle chain<br></li><li>dependency injection support*<br></li><li>context support JSON/JSONP/XML/HTML response<br></li><li>centralized HTTP error handling<br></li><li>centralized log handling<br></li><li>whichever template engine support(emplement baa.Renderer)<br></li></ul><p><b>Middleware</b></p><ul><li><a href=\"http://www.example.com\" target=\"_blank\"></a><a href=\"https://github.com/baa-middleware/gzip\" target=\"_self\">gzip</a><br></li><li><a href=\"https://github.com/baa-middleware/logger\" target=\"_blank\">logger</a><br></li><li><a href=\"https://github.com/baa-middleware/recovery\" target=\"_blank\">recovery</a><br></li><li><a href=\"https://github.com/baa-middleware/session\" target=\"_blank\">session</a><br></li></ul><h3><b>Route Test</b></h3><p>Based on&nbsp;<a href=\"https://github.com/safeie/go-http-routing-benchmark\">go-http-routing-benchmark</a>, Feb 27, 2016.</p><p><a href=\"http://developer.github.com/v3\">GitHub API</a></p><p>Baa route test is very close to Echo.</p><pre><code>BenchmarkBaa_GithubAll                 30000         50984 ns/op           0 B/op          0 allocs/op\r\nBenchmarkBeego_GithubAll                3000        478556 ns/op        6496 B/op        203 allocs/op\r\nBenchmarkEcho_GithubAll                30000         47121 ns/op           0 B/op          0 allocs/op\r\nBenchmarkGin_GithubAll                 30000         41004 ns/op           0 B/op          0 allocs/op\r\nBenchmarkGocraftWeb_GithubAll           3000        450709 ns/op      131656 B/op       1686 allocs/op\r\nBenchmarkGorillaMux_GithubAll            200       6591485 ns/op      154880 B/op       2469 allocs/op\r\nBenchmarkMacaron_GithubAll              2000        679559 ns/op      201140 B/op       1803 allocs/op\r\nBenchmarkMartini_GithubAll               300       5680389 ns/op      228216 B/op       2483 allocs/op\r\nBenchmarkRevel_GithubAll                1000       1413894 ns/op      337424 B/op       5512 allocs/op\r\n</code></pre><h3><b>HTTP Test</b></h3><h4>Code</h4><p>Baa:</p><pre><code>package main\r\n\r\nimport (\r\n    \"github.com/baa-middleware/logger\"\r\n    \"github.com/baa-middleware/recovery\"\r\n    \"gopkg.in/baa.v1\"\r\n)\r\n\r\nfunc hello(c *baa.Context) {\r\n    c.String(200, \"Hello, World!\\n\")\r\n}\r\n\r\nfunc main() {\r\n    b := baa.New()\r\n    b.Use(logger.Logger())\r\n    b.Use(recovery.Recovery())\r\n\r\n    b.Get(\"/\", hello)\r\n\r\n    b.Run(\":8001\")\r\n}\r\n</code></pre><p>Echo:</p><pre><code>package main\r\n\r\nimport (\r\n    \"github.com/labstack/echo\"\r\n    mw \"github.com/labstack/echo/middleware\"\r\n)\r\n\r\n// Handler\r\nfunc hello(c *echo.Context) error {\r\n    return c.String(200, \"Hello, World!\\n\")\r\n}\r\n\r\nfunc main() {\r\n    // Echo instance\r\n    e := echo.New()\r\n\r\n    // Middleware\r\n    e.Use(mw.Logger())\r\n\r\n    // Routes\r\n    e.Get(\"/\", hello)\r\n\r\n    // Start server\r\n    e.Run(\":8001\")\r\n}\r\n</code></pre><h4>Result:</h4><p>Baa http test is almost better than Echo.</p><p>Baa:</p><pre><code>$ wrk -t 10 -c 100 -d 30 http://127.0.0.1:8001/\r\nRunning 30s test @ http://127.0.0.1:8001/\r\n  10 threads and 100 connections\r\n  Thread Stats   Avg      Stdev     Max   +/- Stdev\r\n    Latency     1.92ms    1.43ms  55.26ms   90.86%\r\n    Req/Sec     5.46k   257.26     6.08k    88.30%\r\n  1629324 requests in 30.00s, 203.55MB read\r\nRequests/sec:  54304.14\r\nTransfer/sec:      6.78MB\r\n</code></pre><p>Echo:</p><pre><code>$ wrk -t 10 -c 100 -d 30 http://127.0.0.1:8001/\r\nRunning 30s test @ http://127.0.0.1:8001/\r\n  10 threads and 100 connections\r\n  Thread Stats   Avg      Stdev     Max   +/- Stdev\r\n    Latency     2.83ms    3.76ms  98.38ms   90.20%\r\n    Req/Sec     4.79k     0.88k   45.22k    96.27%\r\n  1431144 requests in 30.10s, 178.79MB read\r\nRequests/sec:  47548.11\r\nTransfer/sec:      5.94MB</code></pre>', '1', '0', '2016-12-25 23:59:02', '2016-12-25 23:59:02');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
