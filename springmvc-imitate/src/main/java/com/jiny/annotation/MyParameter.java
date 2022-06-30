package com.jiny.annotation;

import java.lang.annotation.*;

/**
 * @Auther: jiny
 * @CreateDate: 2018/9/16 17:18
 * @Description: 声明参数
 */
@Target(ElementType.PARAMETER)//只能声明在参数上
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface MyParameter {
    String value();
}
