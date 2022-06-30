
## springmvc-imitate
### 介绍
模仿spring mvc实现了控制器注册、扫描和访问。

### 使用方式
- clone 下来项目打包到本地maven
- pom 引用
- web.xml配置控制器
```
<servlet>
   <servlet-name>springmvc</servlet-name>
   <servlet-class>com.jiny.core.MyDispatcherServlet</servlet-class>
   <init-param>
            <param-name>contextConfigLocation</param-name>
            <param-value>classpath:spring-mvc.xml</param-value>
   </init-param>
   <load-on-startup>1</load-on-startup>
</servlet>
<servlet-mapping>
   <servlet-name>springmvc</servlet-name>
   <url-pattern>/*</url-pattern>
</servlet-mapping>
```
- spring-mvc.xml编写
```
<?xml version="1.0" encoding="UTF-8"?>
<beans>
    <component-scan base-package="com.test.controller.*"></component-scan>
</beans>
```
- 注解的使用在原有的spring mvc的注解前面加上My. eg : @MyController