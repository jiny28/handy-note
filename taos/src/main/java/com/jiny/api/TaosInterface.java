package com.jiny.api;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/5 10:35
 * @Description:
 */
public interface TaosInterface {




    Boolean updateSubTableTag(String subTable, String tagName, String tagValue);



    /**
    * @Author:         jiny
    * @CreateDate:     2022/7/8 10:49
    * @Description:    由超级表做模板插入子表，需要传入tag值，子表自动创建。
    */
    Boolean insertSubTable();

}
