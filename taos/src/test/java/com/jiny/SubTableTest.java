package com.jiny;

import com.jiny.api.QueryInterface;
import com.jiny.api.SubTableInterface;
import com.jiny.entity.*;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.*;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.Is.is;

@SpringBootTest
class SubTableTest {

    @Autowired
    private SubTableInterface subTableInterface;


    @Test
    void create() {
        subTableInterface.create(new SubTableMeta("", "meters", "d001", new ArrayList<TagValue>() {{
            add(new TagValue<>("location", "123"));
            add(new TagValue<>("groupId", 123));
        }}, null));
    }

    // 存在sql语句过长问题
    @Test
    void insert() throws ParseException {
        List<SubTableValue> subTableValues = new ArrayList<>();
        Random random = new Random();
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        String startTime = "2022-07-12 10:00:00";
        long end = new Date().getTime();
        long start = sdf.parse(startTime).getTime();
        for (int i = 0; i < 1; i++) {
            SubTableValue subTableValue = new SubTableValue();
            subTableValue.setDatabase("");
            subTableValue.setName("d001");
            List<RowValue> rowValues = new ArrayList<>();
            int j = 1;
            while (true) {
                List<FieldValue> fieldValues = new ArrayList<>();
                fieldValues.add(new FieldValue("ts", start));
                fieldValues.add(new FieldValue("current", random.nextFloat()));
                fieldValues.add(new FieldValue("voltage", j++));
                fieldValues.add(new FieldValue("phase", random.nextFloat()));
                RowValue rowValue = new RowValue(fieldValues);
                rowValues.add(rowValue);
                start += 1000;
                if (start >= end) {
                    break;
                }
            }
            subTableValue.setValues(rowValues);
            subTableValues.add(subTableValue);
        }
        int insert = subTableInterface.insert(subTableValues);
        assertThat(insert, is(subTableValues.size() * subTableValues.get(0).getValues().size()));
    }

    @Test
    void insertAutoCreateTable() throws ParseException {
        List<SubTableValue> subTableValues = new ArrayList<>();
        Random random = new Random();
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        String startTime = "2022-07-10 10:00:00";
        long start = sdf.parse(startTime).getTime();
        // 100个电表，每块电表200个点位
        for (int i = 0; i < 1; i++) {
            SubTableValue subTableValue = new SubTableValue();
            subTableValue.setDatabase("");
            subTableValue.setName("d00" + i);
            subTableValue.setSuperTable("meters");
            List<TagValue> tags = new ArrayList<>();
            tags.add(new TagValue<>("location", "d00" + i));
            tags.add(new TagValue<>("groupId", 1 + i));
            subTableValue.setTags(tags);
            List<RowValue> rowValues = new ArrayList<>();
            for (int j = 0; j < 100; j++) {//一次插入多少数据
                List<FieldValue> fieldValues = new ArrayList<>();
                fieldValues.add(new FieldValue("ts", start));
                for (int a = 0; a < 200; a++) {
                    String fieldName = "field" + a;
                    float v = random.nextFloat();
                    int f = random.nextInt();
                    FieldValue fieldValue;
                    if (a % 2 == 0) {
                        fieldValue = new FieldValue(fieldName, v);
                    } else {
                        fieldValue = new FieldValue(fieldName, f);
                    }
                    fieldValues.add(fieldValue);
                }
                rowValues.add(new RowValue(fieldValues));
                start += 1000;
            }
            subTableValue.setValues(rowValues);
            subTableValues.add(subTableValue);
        }
        long s = System.currentTimeMillis();
        int insert = subTableInterface.insertAutoCreateTable(subTableValues);
        long e = System.currentTimeMillis();
        System.out.println("插入用时:" + (e - s));
        /*
        * 跑sql最长限制执行时间：722; 一张表，19300条数据，4个字段
        * 100张表，每张表200个字段，每张表插入1条数据，共计100条数据， 耗时1988
        * 1张表，每张表200个字段，每张插入100条数据，共计100条数据， 耗时485
        *
        * */
        assertThat(insert, is(subTableValues.size() * subTableValues.get(0).getValues().size()));

    }


    @Test
    void updateSubTableTag() {
        // 客户端需要配置FQDN，host设置映射到容器id或者是容器hostname
        subTableInterface.updateSubTableTag("", "d001", "groupid", "444");
    }

    @Autowired
    private QueryInterface queryInterface;


    @Test
    void select() {
        List<LinkedHashMap<String, String>> list = queryInterface.executeQuery("select * from meters");
        System.out.println(list.size());
    }

}
