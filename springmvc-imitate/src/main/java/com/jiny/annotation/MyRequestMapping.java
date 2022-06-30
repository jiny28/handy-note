package com.jiny.annotation;

import java.lang.annotation.*;

/**
* @Author:         jiny
* @CreateDate:     2019/7/8 10:14
* @Description:    装载class method
*/
//作用于类，接口，方法上
@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)//加载到vm
@Documented
public @interface MyRequestMapping {
    //url 必填
    String value();
}
