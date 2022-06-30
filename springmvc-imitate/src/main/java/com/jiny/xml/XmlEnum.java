package com.jiny.xml;

/**
 * @Auther: jiny
 * @CreateDate: 2018/9/16 12:14
 * @Description: 用于声明xml文件的配置规则
 */
public enum XmlEnum {
    SCAN_RULE("component-scan", "base-package", "null");

    private String type;
    private String name;
    private String value;
    XmlEnum(String property, String name, String value) {
        this.type  = property;
        this.name  = name;
        this.value = value;
    }

    public String getType() {
        return type;
    }



    public String getName() {
        return name;
    }



    public String getValue() {
        return value;
    }

}
