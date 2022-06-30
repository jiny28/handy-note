package com.jiny.annotation;

import java.lang.annotation.*;

/**
* @Author:         jiny
* @CreateDate:     2019/7/8 10:13
* @Description:    装载class
*/
@Target(ElementType.TYPE)//作用于类和接口上
@Retention(RetentionPolicy.RUNTIME)//加载到vm
@Documented
public @interface MyController {

    //自定义别名
    String value() default "";
}
