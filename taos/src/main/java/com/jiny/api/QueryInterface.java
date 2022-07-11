package com.jiny.api;

import java.util.LinkedHashMap;
import java.util.List;

/**
 * @Auther: jiny
 * @CreateDate: 2022/7/11 16:15
 * @Description:
 */
public interface QueryInterface {


    List<LinkedHashMap<String, String>> executeQuery(String sql);

}
