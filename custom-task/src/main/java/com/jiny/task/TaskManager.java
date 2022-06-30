package com.jiny.task;


import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;


/**
 * @Author: jiny
 * @CreateDate: 2019/6/18 13:40
 * @Description: 定时任务分组管理类
 */
public class TaskManager {

    //所有分组
    private static Map<String, GroupTask> group;

    static {
        group = new ConcurrentHashMap<>();
    }

    /**
     * @Author: jiny
     * @CreateDate: 2019/6/18 10:58
     * @Description: 添加定时任务分组
     */
    public static void addGroup(String name, GroupTask tg) {
        if (group.containsKey(name)) {
            return;
        }
        group.put(name, tg);
    }

    /**
     * @Author: jiny
     * @CreateDate: 2019/6/18 11:01
     * @Description: 杀掉定时任务分组
     */
    public static Boolean killGroup(String name) {
        if (group.containsKey(name)) {
            group.get(name).setRunFlag(false);
            group.remove(name);
            return true;
        }
        return false;
    }

    /**
    * @Author:         jiny
    * @CreateDate:     2019/7/19 15:24
    * @Description:    停止分组中的一个任务
    */
    public static Boolean killJobByGroup(String groupName,String jobName){
        GroupTask groupTask = group.get(groupName);
        if (groupTask == null) {
            return false;
        }
        if (!groupTask.getAllTask().containsKey(jobName)) {
            return false;
        }
        JobTask jobTask = groupTask.getAllTask().get(jobName);
        jobTask.setRunFlag(false);
        jobTask.setRuning(false);
        jobTask.setGuarding(false);
        return true;
    }

    /**
    * @Author:         jiny
    * @CreateDate:     2019/7/19 15:36
    * @Description:    开始分组中的所有任务
    */
    public static Boolean startAllJobByGroup(String groupName){
        GroupTask groupTask = group.get(groupName);
        if (groupTask == null) {
            return false;
        }
        for (Map.Entry<String, JobTask> entry : groupTask.getAllTask().entrySet()) {
            startJobByGroup(groupName, entry.getKey());
        }
        return true;
    }

    /**
    * @Author:         jiny
    * @CreateDate:     2019/7/19 15:34
    * @Description:    开始分组中的一个任务
    */
    public static Boolean startJobByGroup(String groupName,String jobName){
        GroupTask groupTask = group.get(groupName);
        if (groupTask == null) {
            return false;
        }
        if (!groupTask.getAllTask().containsKey(jobName)) {
            return false;
        }
        JobTask jobTask = groupTask.getAllTask().get(jobName);
        //通过守护线程是否开启判断该定时任务是否开启
        if (!jobTask.getRunFlag()) {
            jobTask.startTask();
        }
        return true;
    }
}
