## custom-task
### 介绍
定时任务调动中心。每开启一个定时任务,程序会开启两个线程,一个线程用于执行任务,另外一个用于守护任务线程,当任务线程挂掉后，守护线程会替代它。
### 使用方式
- 定义任意类(job)实现ITaskHandle接口
- 创建分组
```
GroupTask groupTask = new GroupTask("group");//任务组名
```
- 分组添加job
```
ITaskHandle test = new Test();//任务实例
groupTask.createTask("test", 1000, test);//为该组创建一个任务，单位为毫秒
```
- 管理类操作分组
```
TaskManager.startAllJobByGroup("group");//开启该分组的所有任务
TaskManager.startJobByGroup("group","test");//开启该分组的指定任务
TaskManager.killJobByGroup("group","test");//停止该分组的所有任务
TaskManager.killGroup("group");//kill该分组
```