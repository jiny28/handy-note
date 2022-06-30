package com.jiny.xml;

import java.util.List;

/**
 * @Auther: jiny
 * @Description: 用于解析xml
 */
public interface XmlApplication {

    /**
    * @Author:         jiny
    * @CreateDate:     2018/9/16 12:06
    * @Description:    根据指定的xml，获得注解扫描的bean容器
    */
    List<String> getComponentList(String contextConfigLocation);
}
