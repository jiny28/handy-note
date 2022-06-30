package com.jiny.task;



/**
* @Author:         jiny
* @CreateDate:     2019/6/18 13:42
* @Description:    每个定时任务的守护任务
*/
public class GuardTask implements Runnable {

    //所需守护的任务
    private com.jiny.task.JobTask jobTask;

    public GuardTask(com.jiny.task.JobTask jobTask) {
        this.jobTask = jobTask;
    }

    public void run() {
        executeGuard();
    }

    /**
    * @Author:         jiny
    * @CreateDate:     2019/6/18 13:45
    * @Description:    守护线程
    */
    private void executeGuard() {
        getJobTask().setGuarding(true);
        while (getJobTask().getRunFlag()) {
            //当任务没在运行的时候才开启线程
            if (!getJobTask().getRuning()) {
                Thread thread = new Thread(getJobTask(), getJobTask().getTaskName());
                thread.start();
                try {
                    thread.join();
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
            getJobTask().sleep(getJobTask().getPeriod() * 10);
        }
        //守护线程结束
        getJobTask().setGuarding(false);
    }


    public com.jiny.task.JobTask getJobTask() {
        return jobTask;
    }

    public void setJobTask(com.jiny.task.JobTask jobTask) {
        this.jobTask = jobTask;
    }
}
