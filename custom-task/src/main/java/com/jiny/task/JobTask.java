package com.jiny.task;


/**
 * @Auther: jiny
 * @CreateDate: 2019/6/18 11:26
 * @Description: 定时任务
 */
public class JobTask implements Runnable {

    //定时任务名称
    private String taskName;
    //周期
    private Integer period;
    //定时任务执行类
    private ITaskHandle iTask;
    //属于哪个定时任务分组
    private GroupTask groupTask;
    //是否开启了守护线程,默认未开启
    private boolean isGuarding;
    //线程是否运行的标识
    private Boolean runFlag;
    //该定时任务是否正在运行中
    private Boolean isRuning;

    public JobTask(String taskName, Integer period, ITaskHandle iTask, GroupTask groupTask) {
        this.taskName = taskName;
        this.period = period;
        this.iTask = iTask;
        this.groupTask = groupTask;
    }

    {
        this.isGuarding = false;
        this.runFlag = false;
        this.isRuning = false;
    }


    /**
     * @Author: jiny
     * @CreateDate: 2019/6/18 12:49
     * @Description: 开启定时任务
     */
    public void startTask() {
        //如果守护线程开启或者是运行状态，则return
        if (isGuarding() || getRuning()) {
            return;
        }
        //标识赋值为true
        setRunFlag(true);
        //开启守护线程
        new Thread(new GuardTask(this), getTaskName() + "Guard").start();
    }

    public void run() {
        doTask();
    }


    /**
     * @Author: jiny
     * @CreateDate: 2019/6/18 13:03
     * @Description: 线程所执行的方法
     */
    private void doTask() {
        //线程开始，则把任务的运行状态换为true
        setRuning(true);
        while (getRunFlag() && groupTask.isRunFlag()) {
            Long now = System.nanoTime() / 1000000;
            if (getiTask() != null) {
                //执行业务方法
                getiTask().hTaskRun();
            }
            //计算业务方法执行了多长时间
            Long timeSpan = System.nanoTime() / 1000000 - now;
            //使用周期减去时间得出sleep多少毫秒
            Long sleep = getPeriod() - timeSpan;
            if (sleep < 10) {
                sleep = 10l;
            }
            sleep(sleep.intValue());
        }
        //当runFlag结束时，同时也得把运行状态修改为false
        setRuning(false);
        setRunFlag(false);
    }


    /**
    * @Author:         jiny
    * @CreateDate:     2019/6/18 15:10
    * @Description:    sleep多少毫秒
    */
    public void sleep(int sleep) {
        int unitTime = 1000;
        int multiple = sleep / unitTime;
        try {
            for (int i = 0; i < multiple; i++) {
                Thread.sleep(unitTime);
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    /**
     * @Author: jiny
     * @CreateDate: 2019/6/18 14:02
     * @Description: get、set方法
     */
    public String getTaskName() {
        return taskName;
    }

    public void setTaskName(String taskName) {
        this.taskName = taskName;
    }

    public Integer getPeriod() {
        return period;
    }

    public void setPeriod(Integer period) {
        this.period = period;
    }

    public ITaskHandle getiTask() {
        return iTask;
    }

    public void setiTask(ITaskHandle iTask) {
        this.iTask = iTask;
    }

    public GroupTask getGroupTask() {
        return groupTask;
    }

    public void setGroupTask(GroupTask groupTask) {
        this.groupTask = groupTask;
    }

    public boolean isGuarding() {
        return isGuarding;
    }

    public void setGuarding(boolean guarding) {
        isGuarding = guarding;
    }

    public Boolean getRunFlag() {
        return runFlag;
    }

    public void setRunFlag(Boolean runFlag) {
        this.runFlag = runFlag;
    }

    public Boolean getRuning() {
        return isRuning;
    }

    public void setRuning(Boolean runing) {
        isRuning = runing;
    }



}
