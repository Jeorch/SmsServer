# 配置文件问题
 
由于在STS的过程中，使用Redis缓存了ali STS Token
在处理的过程中，我以前使用了死路径读取Redis配置文件，照成了现阶段数据配置文件重复的问题。
现在先在resource 和 Resources 两个文件夹中都存在 redisconfig.json文件

# 现阶段保留两个文件，以后再修改
