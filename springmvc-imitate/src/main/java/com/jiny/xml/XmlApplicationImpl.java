package com.jiny.xml;

import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.io.SAXReader;

import java.io.File;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;

/**
 * @Auther: jiny
 * @Description: 解析xml
 */
public class XmlApplicationImpl implements XmlApplication {
    @Override
    public List<String> getComponentList(String contextConfigLocation) {
        List<String> componentList = new ArrayList<>();
        List<Element> elements = getElements(contextConfigLocation);
        if (elements == null) {
            throw new RuntimeException("xml is null");
        }
        for (Element element : elements) {
            //等于scan的name
            if (element.getName().equals(XmlEnum.SCAN_RULE.getType())) {
                //取到包名
                String packageName = element.attributeValue(XmlEnum.SCAN_RULE.getName());
                componentList.addAll(scanPackage(packageName));
            }
        }
        return componentList;
    }
    
    /**
    * @Author:         jiny
    * @CreateDate:     2018/9/16 12:12
    * @Description:    获取配置文件的所有子元素
    */
    public List<Element> getElements(String contextConfigLocation){
        SAXReader reader = new SAXReader();
        Document read = null;
        try {
            //读取配置文件
            read = reader.read(contextConfigLocation);
        } catch (DocumentException e) {
            e.printStackTrace();
            return null;
        }
        Element rootElement = read.getRootElement();
        if (rootElement == null) {
            return null;
        }
        return rootElement.elements();
    }
    /**
    * @Author:         jiny
    * @CreateDate:     2018/9/16 12:29
    * @Description:    根据要扫描的包名，返回下面所有的类
    */
    public List<String> scanPackage(String packageName){
        if (packageName == null || "".equals(packageName)) {
            throw new RuntimeException(XmlEnum.SCAN_RULE.getType() + "is error");
        }
        List<String> classNames = new ArrayList<>();
        getClassName(packageName, classNames);
        return classNames;
    }

    /**
    * @Author:         jiny
    * @CreateDate:     2018/9/16 13:36
    * @Description:    递归扫描地址下面的classname
    */
    private void getClassName(String packageName, List<String> classNames) {
        //把所有的.*删除掉 并且把.替换成/
        packageName = packageName.replace(".*", "");
        String newName = "/" + packageName.replaceAll("\\.", "/");
        URL url = this.getClass().getClassLoader().getResource(newName);
        File dir = new File(url.getFile());
        for (File file : dir.listFiles()) {
            if(file.isDirectory()){
                //递归读取包
                getClassName(packageName+"."+file.getName(),classNames);
            }else if (file.getName().endsWith(".class")){
                //文件名为.class才添加
                String className =packageName +"." +file.getName().replace(".class", "");
                classNames.add(className);
            }
        }
    }
}
