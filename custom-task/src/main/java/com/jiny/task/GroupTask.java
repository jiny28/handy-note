package com.jiny.task;


import java.util.HashMap;
import java.util.Map;


/**
 * @Author: jiny
 * @CreateDate: 2019/6/18 10:56
 * @Description: 定时任务分组
 */
public class GroupTask {

    //分组名称
    private String groupName;
    //所有定时任务
    private Map<String, JobTask> allTask;
    //运行标识
    private boolean runFlag;

    {
        this.allTask = new HashMap<String, JobTask>();
        this.runFlag = true;
    }

    public GroupTask(String groupName) {
        this.groupName = groupName;
        //把实例添加进管理池
        TaskManager.addGroup(groupName, this);
    }



    /**
    * @Author:         jiny
    * @CreateDate:     2019/6/18 11:04
    * @Description:    创建定时任务
     * @param taskName 定时任务名称
     * @param period 定时任务执行周期,单位为毫秒
     * @param executeClass 实现适配器的执行类
    */
    public JobTask createTask(String taskName, int period, ITaskHandle executeClass) {
        //1. 判断该线程分组是否启用,如果被杀掉，则不能添加定时任务
        if (!isRunFlag()) {
            return null;
        }
        //2. 判断该定时任务是否存在，如果存在，则判断该定时任务是否在运行，如果没运行，还是可以加入进去，
        //加入进去后需要重新启动才能启动
        if (getAllTask().containsKey(taskName) && getAllTask().get(taskName).getRuning()) {
            return null;
        }
        //3. 以上条件都通过后再创建定时任务
        // 传group实例进去的目的：
        //1 . 全局控制所有定时任务
        JobTask jobTask = new JobTask(taskName, period, executeClass, this);
        allTask.put(taskName, jobTask);
        return jobTask;
    }



    public String getGroupName() {
        return groupName;
    }

    public void setGroupName(String groupName) {
        this.groupName = groupName;
    }

    public Map<String, JobTask> getAllTask() {
        return allTask;
    }

    public void setAllTask(Map<String, JobTask> allTask) {
        this.allTask = allTask;
    }

    public boolean isRunFlag() {
        return runFlag;
    }

    public void setRunFlag(boolean runFlag) {
        this.runFlag = runFlag;
    }



}
