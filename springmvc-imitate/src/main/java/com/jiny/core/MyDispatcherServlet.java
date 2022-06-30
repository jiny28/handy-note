package com.jiny.core;

import com.jiny.annotation.MyController;
import com.jiny.annotation.MyParameter;
import com.jiny.annotation.MyRequestMapping;
import com.jiny.xml.XmlApplication;
import com.jiny.xml.XmlApplicationImpl;

import javax.servlet.ServletConfig;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.lang.reflect.Constructor;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.util.*;

/**
* @Author:         jiny
* @CreateDate:     2019/7/8 11:12
* @Description:    控制器
*/
public class MyDispatcherServlet extends HttpServlet {

    //存储控制器所控制的所有类与实例
    private Map<String, Object> ioc = new HashMap<>();

    //适配器,存储所有被控制的url-method
    private Map<String, Method> handlerMapping = new HashMap<>();

    //适配器,存储所有被控制的url-controller
    private Map<String, Object> controllerMap = new HashMap<>();

    @Override
    public void init(ServletConfig config) throws ServletException {
        XmlApplication xmlApplication = new XmlApplicationImpl();

        //扫描配置文件里面所定义的包下面的所有类
        String contextConfigLocation = config.getInitParameter("contextConfigLocation");
        List<String> classNames = xmlApplication.getComponentList(contextConfigLocation);

        //找出所有类里面带有@controller注解的类放入ioc容器中
        loadIoc(classNames);

        //加载handlerMapping
        loadHandlerMapping();
    }

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        doPost(req, resp);
    }

    @Override
    protected void doPost(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        try {
            resp.setHeader("Content-type", "text/html;charset=UTF-8");
            doDispatch(req, resp);
        } catch (Exception e) {
            e.printStackTrace();
            resp.getWriter().write("500 error");
        }
    }

    /**
     * @Author: jiny
     * @CreateDate: 2018/9/16 15:27
     * @Description: 执行请求
     */
    private void doDispatch(HttpServletRequest req, HttpServletResponse resp) throws Exception {
        if (handlerMapping.isEmpty()) {
            return;
        }
        //获取url取handelmapping中匹配
        String requestURI = req.getRequestURI();
        String contextPath = req.getContextPath();
        requestURI = requestURI.replace(contextPath, "");
        if (!handlerMapping.containsKey(requestURI)) {
            resp.getWriter().write("404 not fount");
            return;
        }
        Method method = handlerMapping.get(requestURI);
        //获取方法的所有参数
        Parameter[] parameters = method.getParameters();
        //封装参数
        List<Object> obj = new LinkedList<>();
        //获取请求的所有参数
        Map<String, String[]> parameterMap = req.getParameterMap();
        for (int i = 0; i < parameters.length; i++) {
            Class<?> type = parameters[i].getType();
            if (type.getSimpleName().equals("HttpServletRequest")) {
                obj.add(i, req);
                continue;
            }
            if (type.getSimpleName().equals("HttpServletResponse")) {
                obj.add(i, resp);
                continue;
            }
            if (parameters[i].isAnnotationPresent(MyParameter.class)) {
                //获取到参数名字
                String paraName = parameters[i].getAnnotation(MyParameter.class).value();
                String[] value = parameterMap.get(paraName);
                //没传参数就push null
                if (value == null || value.length == 0) {
                    obj.add(i, null);
                    continue;
                }
                StringBuilder sb = new StringBuilder();
                for (String item : value) {
                    sb.append(item);
                }
                //将String转为参数的类型
                Constructor<?> constructor = type.getConstructor(String.class);
                constructor.setAccessible(true);
                Object o = constructor.newInstance(sb.toString());
                obj.add(i, o);
            } else {
                //因为jdk8以下通过反射不能取到参数名字，所以此处只能用注解形式实现
                obj.add(i, null);
            }
        }
        method.invoke(controllerMap.get(requestURI), obj.toArray());
    }


    /**
     * @Author: jiny
     * @CreateDate: 2018/9/16 14:24
     * @Description: 加载handlerMapping
     */
    private void loadHandlerMapping() {
        if (ioc.isEmpty()) {
            return;
        }
        for (Map.Entry<String, Object> entry : ioc.entrySet()) {
            Class<?> clazz = entry.getValue().getClass();
            String baseUrl = "";
            //先判断该类上面有没有@requestmapping
            if (clazz.isAnnotationPresent(MyRequestMapping.class)) {
                baseUrl = clazz.getAnnotation(MyRequestMapping.class).value();
            }
            Method[] methods = clazz.getDeclaredMethods();
            for (Method method : methods) {
                //如果该方法存在注解
                if (method.isAnnotationPresent(MyRequestMapping.class)) {
                    //拼接url
                    String url = baseUrl + method.getAnnotation(MyRequestMapping.class).value();
                    handlerMapping.put(url, method);
                    controllerMap.put(url, entry.getValue());
                }
            }
        }
    }

    /**
     * @Author: jiny
     * @CreateDate: 2018/9/16 14:14
     * @Description: 加载ioc
     */
    private void loadIoc(List<String> classNames) {
        if (classNames.isEmpty()) {
            return;
        }
        for (String name : classNames) {
            try {
                Class<?> clazz = Class.forName(name);
                //判断该clazz是否有某个注解
                if (clazz.isAnnotationPresent(MyController.class)) {
                    //类名字的手写字母小写存储
                    String str = strToLowerFirst(clazz.getSimpleName());
                    ioc.put(str, clazz.newInstance());
                }
            } catch (ClassNotFoundException e) {
                e.printStackTrace();
            } catch (IllegalAccessException e) {
                e.printStackTrace();
            } catch (InstantiationException e) {
                e.printStackTrace();
            }
        }
    }


    /**
     * @Author: jiny
     * @CreateDate: 2018/9/16 14:23
     * @Description: 字符串首字母大写
     */
    private String strToLowerFirst(String name) {
        char[] charArray = name.toCharArray();
        charArray[0] += 32;
        return String.valueOf(charArray);
    }


}
