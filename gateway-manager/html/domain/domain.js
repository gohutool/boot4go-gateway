function load() {
    jquery(function () {
        var productDg = {
            type: 'datagrid',
            id: 'productDg'
        };

        $("#productDg").iDatagrid({
            url: V3_API_URL + '/domain/list',
            idField: 'ID',
            sortName: 'COMPANYNO',
            sortOrder: "asc",
            frozenColumns:[[
                // {field: 'id', title: '', checkbox: true},
                {field: 'id', title: 'ID', sortable: false,
                    formatter:$.iGrid.templateformatter('{id}'),
                    width: 350},
            ]],
            columns: [[
                {
                    field: 'domain_name',
                    title: '名称',
                    sortable: true,
                    width: 180,
                    formatter:$.iGrid.tooltipformatter()
                },
                {field: 'domain_url', title: '域名', sortable: false, width: 240,
                    formatter: function(value, row, index){
                        return '<a onclick=\"openEditDomainDiag(\'' + row.id + '\');return false\" class=\'cubeui-link\' href=\'#\'>'+value+'</a>'
                    }},
                {field: 'set_time', title: '创建日期', sortable: false, width: 160},
                {field: 'op', title: '操作', sortable: false, width: 380, formatter:operateFormatter}
            ]],
            loadFilter: function (data) {
                // add checkedRows test data, just for debug
                var d = $.fn.iDatagrid.defaults.loadFilter.call(this, data);
                //d.checkedRows = [999, 22];
                return d;
            }
        });
    });
}

function operateFormatter(value, row, index) {
    var htmlstr = "";

    htmlstr += '<button class="layui-btn cubeui-btn-ivory layui-btn-xs" onclick="openEditDiag(\'' + row.id + '\')">修改信息</button>';
    htmlstr += '<button class="layui-btn layui-btn-danger layui-btn-xs" onclick="openEditDomainDiag(\'' + row.id + '\')">修改域名</button>';
    htmlstr += '<button class="layui-btn layui-btn-danger layui-btn-xs" onclick="openPathDiag(\'' + row.id + '\')">路径映射</button>';
    htmlstr += '<button class="cubeui-btn-gray layui-btn layui-btn-xs" onclick="delDomain(\'' + row.id + '\')">删除域名</button>';

    return htmlstr;
}

function openEditDomainDiag(id){

    let _d = function (data, fn){
        $.iDialog.openDialog({
            title: '修改域名',
            minimizable:false,
            content: `
            <div style="margin: 0px;">
            </div>
            <div class="cubeui-fluid" id="company-detail-form">
                <div class="cubeui-row">
                    <label class="cubeui-form-label">域名:</label>
                    <div class="cubeui-input-block">
                        <input type="text" data-toggle="cubeui-textbox" id="domain_url" name="domain_url"
                               value='{{:domain_url}}' data-options="required:true,prompt:'域名+端口，必须填写'">
                    </div>
                </div>
                
                <div class="cubeui-row">
                
                    <div class="cubeui-col-sm8">
                        <label class="cubeui-form-label">SSL端口:</label>
                        <div class="cubeui-input-block">
                            <input type="text" disabled data-toggle="cubeui-numberspinner" id="ssl_port" name="ssl_port"
                                   value='{{:ssl_port}}' data-options="required:false,prompt:'SSL端口，默认为9443',
                                   min:10,max:65525,increment:1">
                        </div>
                    </div>
                    
                    <div class="cubeui-col-sm4">
                        <label class="cubeui-form-label">代理SSL:</label>
                        <div class="cubeui-input-block">
                            <input type="text" data-toggle="cubeui-switchbutton" id="ssl_on" name="ssl_on"
                                   value='{{:ssl_on}}' data-options="value:1,width:'50px',required:false,onText:'',offText:''">
                        </div>
                    </div>
                    
                </div>
                
            </div>
            `,
            width: 600,
            height: 300,
            render:function(opts, handler){
                let d = this;
                console.log("Open dialog")

                handler.render(data)

                $("#ssl_on").switchbutton('options').onChange = function(checked){
                    if(checked){
                        $("#ssl_port").numberspinner('enable')
                    }else{
                        $("#ssl_port").numberspinner('disable')
                    }
                }

                if (data.ssl_on+'' == 'false') {
                    $("#ssl_on").switchbutton('uncheck')
                }else{
                    $("#ssl_on").switchbutton('check')
                }

                if(fn){
                    fn.call(d, opts)
                }

            },
            buttonsGroup: [{
                text: '保存',
                iconCls: 'fa fa-save',
                btnCls: 'cubeui-btn-blue',
                handler:'ajaxForm',
                requestType:'put',
                reload: [{type:'datagrid', id:'productDg', keepcheck:0}],
                url:V3_API_URL + '/domain/url/'+id
            }]
        });
    }


    $.app.get(V3_API_URL + '/domain/'+id, null, function (data) {
        _d(data.data)
    });
}

function delDomain(id){
    $.app.confirm(null, "确认删除域名"+id+", 域名下的路径映射规则将一并删除?", function (){
        $.app.deleteJson(V3_API_URL + '/domain/'+id, null, function (data) {
            $.app.show(data.msg);
            $('#productDg').datagrid('reload')
        })
    })
}

function openEditDiag(id){

    _d = function(data, fn){
        $.iDialog.openDialog({
            title: '编辑',
            maximized1:true,
            minimizable:false,
            width: 780,
            height: 400,
            render:function(opts, handler){
                let d = this;
                console.log("Open dialog")
                handler.render(data)

                if(fn){
                    fn.call(d, opts)
                }
            },
            href: contextpath + '/domain/domain-detail.html?id='+id,
            buttonsGroup: [{
                text: '保存',
                iconCls: 'fa fa-save',
                reload: [{type:'datagrid', id:'productDg', keepcheck:0}],
                btnCls: 'cubeui-btn-blue',
                handler:'ajaxForm',
                requestType:'put',
                beforeAjax:function(o){
                    o.ajaxData = $.extends.json.param2json(o.ajaxData);
                    o.ajaxData.black_ips = $("#black_ips").tagbox('getValues')
                },
                postJson:true,
                url:V3_API_URL + '/domain/'+(id||'add')
            }]
        });
    }

    if(id!=null){
        $.app.get(V3_API_URL + '/domain/'+id, null, function (data) {
            _d(data.data, function (){
                $(this).dialog('setTitle', '域名：' + id)

                $("#domain_url").textbox('readonly')
            })
        });
    }  else{
        _d({id:'add',ssl_port:9443}, function (){
            $(this).dialog('setTitle', '新增域名')

            $('#domainIdDiv').hide()
            $("#domain_url").textbox('readonly', false)
            $("#domain_name").textbox('readonly', false)
            //initContactDg({})
        })
    }
}

function initContactDg(data){

    $("#contactDg").iDatagrid({
        pagination:false,
        idField:'ID',
        remoteSort:false,
        sortOrder:"asc",
        data:{rows:(data&&data.data&&data.data.targets)?data.data.targets:[]},
        singleSelect:true,
        onDblClickRow:$.iGrid.EdatagridHandle.rowEditRowEventHanle(function(idx, row, colopts){
        }),
        columns: [[
            {field: 'ID', title: '', checkbox: true},
            {field: 'pointer_type', title: '目标类型',width:450, editor:{type:'combobox',
                    options:{required:true, prompt:"目标类型(单机/集群/Etcd/Nacos/ZK)，必须填写",
                        readonly1:true,valueField:'code',value:'Node',
                        data:[
                            {code:'Node', text:'单机'},
                            {code:'Cluster', text:'集群'},
                            {code:'Etcd', text:'Etcd服务治理'},
                            {code:'Nacos', text:'Nacos服务治理'},
                            {code:'ZK', text:'ZooKeeper服务治理'}
                        ]}
                },
                sortable: true,
                formatter:$.iGrid.evalconditionformatter([
                        'rowData[\'pointer_type\']=="Node"', 'rowData[\'pointer_type\']=="ZK"',
                        'rowData[\'pointer_type\']=="Cluster"', 'rowData[\'pointer_type\']=="Etcd"',
                        'rowData[\'pointer_type\']=="Nacos"'],
                    ['单机', 'ZooKeeper服务治理', '集群', 'Etcd服务治理', 'ZooKeeper服务治理'], '单机')
            },
            {field: 'pointer', title: '目标地址',width:350, editor:{type:'textbox',options:{required:true, prompt:"目标地址(域名/IP+Port)，必须填写"}},
                sortable: true, formatter:$.iGrid.tooltipformatter()},
            {field: 'weight', title: '权重',width:150, editor:{type:'numberspinner',options:{required:true,prompt:"目标地址权重",value:0, min:0,increment:10}}, sortable: true,
                formatter:$.iGrid.tooltipformatter()},
        ]]
    });
}

function openPathDiag(id){

    $.iDialog.openDialog({
        title: '编辑',
        maximized1:true,
        minimizable:false,
        width: 1350,
        height: 650,
        render:function(opts, handler){
            let d = this;
            console.log("Open dialog")

            $.app.get(V3_API_URL + '/path/'+id, null, function (data) {
                data.data.domainId = id;
                handler.render(data.data)

                $(d).dialog('setTitle', '域名路径映射：' + id)

                initPathList(data)
            });
        },
        href: contextpath + '/domain/path.html?id='+id
    });

}

function refreshPathList(domainId){
    $.app.get(V3_API_URL + '/path/'+domainId, null, function (data) {
        $("#pathDg").datagrid("loadData", data.data.list);
    });
}

function initPathList(data){

    $("#pathDg").iDatagrid({
        idField:'ID',
        remoteSort:false,
        sortOrder:"asc",
        data:ConvertResult2Data(data),
        frozenColumns:[[
            {field: 'id', title: 'ID',  width:300,
                sortable: true, formatter:$.iGrid.tooltipformatter()},
            {field: 'req_name', title: '映射描述',  width:240,
                sortable: true, formatter:$.iGrid.tooltipformatter()},
            {field: 'op', title: '操作', sortable: false, width: 100, formatter:wrapPathOperate(data.data.domainId)},
        ]],
        columns: [[
            {field: 'req_path', title: '请求路径',width:120,
                sortable: true, formatter:$.iGrid.tooltipformatter()},
            {field: 'req_method', title: '请求方式',width:80,
                sortable: true, formatter:$.iGrid.tooltipformatter()},
            {field: 'search_path', title: '路径格式',width:120,
                sortable: true, formatter:$.iGrid.tooltipformatter()},
            {field: 'replace_path', title: '转换格式',width:120, sortable: true,
                formatter:$.iGrid.tooltipformatter()},
            {field: 'private_proxy_enabled', title: '私有转发',width:80, sortable: true,
                formatter:$.iGrid.checkedformatter("true", "<span class=\"cube-label cube-label-orange\">是</span>", "")},
            {field: 'circuit_breaker_enabled', title: '开启熔断',width:80, sortable: true,
                formatter:$.iGrid.checkedformatter("true", "<span class=\"cube-label cube-label-orange\">开启</span>", "")},
            {field: 'circuit_breaker_force', title: '强制熔断',width:80, sortable: true,
                formatter:$.iGrid.checkedformatter(true, "<span class=\"cube-label cube-label-orange\">强制</span>", "")},

        ]]
    });

    $("#pathDg").datagrid('getPager').pagination({
        layout:['info']
    });
}

function wrapPathOperate(domainId){
    return function (value, row, index) {
        var htmlstr = "";

        htmlstr += '<button class="layui-btn cubeui-btn-ivory layui-btn-xs" type="button" onclick="openPathDetailDiag(\''+domainId+'\', \'' + row.id + '\')">修改</button>';
        htmlstr += '<button class="cubeui-btn-gray layui-btn layui-btn-xs" type="button" onclick="delPathDetail(\''+domainId+'\', \'' + row.id + '\')">删除</button>';

        return htmlstr;
    }
}

function openPathDetailDiag(domainId, pathId){

    if($.extends.isEmpty(pathId) || pathId === "add" ){
        pathId = "";
    }

    $.iDialog.openDialog({
        title: '编辑',
        id:'pathDetailDlg',
        maximized:true,
        minimizable:false,
        width: 1150,
        height: 600,
        render:function(opts, handler){
            let d = this;
            console.log("Open dialog")

            if(pathId == null || $.extends.isEmpty(pathId) || pathId === "add" ){
                handler.render({})
                $(d).dialog('setTitle', '域名路径映射：添加')

                $("#domain_url").textbox('readonly', false)
                $("#domain_name").textbox('readonly', false)

                $("#rate_limiter_enabled").switchbutton('options').onChange = function(checked){
                    if(checked){
                        $("#rate_limiter_num").numberspinner('enable')
                        $("#rate_limiter_msg").numberspinner('enable')
                    }else{
                        $("#rate_limiter_num").numberspinner('disable')
                        $("#rate_limiter_msg").numberspinner('disable')
                    }
                }

                $("#rate_limiter_enabled").switchbutton('uncheck')

                initContactDg({})
            }else{
                $.app.get(V3_API_URL + '/path/'+domainId+"/"+pathId, null, function (data) {
                    handler.render(data.data)
                    $(d).dialog('setTitle', '域名路径映射：' + pathId)

                    if (data.data.private_proxy_enabled+'' == 'false') {
                        $("#private_proxy_enabled").switchbutton('uncheck')
                    }else{
                        $("#private_proxy_enabled").switchbutton('check')
                    }

                    if (data.data.circuit_breaker_enabled+'' == 'false') {
                        $("#circuit_breaker_enabled").switchbutton('uncheck')
                    }else{
                        $("#circuit_breaker_enabled").switchbutton('check')
                    }

                    if (data.data.circuit_breaker_force+'' == 'false') {
                        $("#circuit_breaker_force").switchbutton('uncheck')
                    }else{
                        $("#circuit_breaker_force").switchbutton('check')
                    }

                    // $("#rate_limiter_enabled").switchbutton('options').onChange = function(checked){
                    //     if(checked){
                    //         $("#rate_limiter_num").numberspinner('enable')
                    //         $("#rate_limiter_msg").numberspinner('enable')
                    //     }else{
                    //         $("#rate_limiter_num").numberspinner('disable')
                    //         $("#rate_limiter_msg").numberspinner('disable')
                    //     }
                    // }

                    $("#domain_url").textbox('readonly')

                    $("#rate_limiter_enabled").switchbutton('options').onChange = function(checked){
                        if(checked){
                            $("#rate_limiter_num").numberspinner('enable')
                            $("#rate_limiter_msg").numberspinner('enable')
                        }else{
                            $("#rate_limiter_num").numberspinner('disable')
                            $("#rate_limiter_msg").numberspinner('disable')
                        }
                    }

                    if (data.data.rate_limiter_enabled == false || data.data.rate_limiter_enabled == 'false') {
                        $("#rate_limiter_enabled").switchbutton('uncheck')
                    }else{
                        $("#rate_limiter_enabled").switchbutton('check')
                    }

                    initContactDg(data)

                });
            }


        },
        href: contextpath + '/domain/path-detail.html?id='+pathId,
        buttonsGroup: [{
            text: '保存',
            iconCls: 'fa fa-save',
            btnCls: 'cubeui-btn-blue',
            handler:'ajaxForm',
            requestType:'put',
            postJson:true,
            beforeAjax:function(o){
                $('#contactDg').datagrid("acceptChanges")
                var rows = $('#contactDg').datagrid('getRows');
                if(rows.length==0){
                    $.app.showerror('请至少添加一个目标地址');
                    return false;
                }
                var d = [];
                $.each(rows, function(idx, row){
                    var o = {};
                    o.pointer = row.pointer;
                    o.pointer_type = row.pointer_type;
                    o.weight = row.weight;
                    d.push(o);
                });

                o.ajaxData = $.extends.json.param2json(o.ajaxData);
                o.ajaxData.targets = d
                //var d = $.param({contacts:$.extends.json.tostring(d)});
                //o.ajaxData = $.extends.json.tostring(o.ajaxData)

                $.app.debug(o.ajaxData);
            },
            reload: [{type:'datagrid', id:'pathDg', keepcheck:0, fn:function(dgObj){
                    console.log("load paths again")
                    refreshPathList(domainId)
                }}],
            url:V3_API_URL + '/path/'+domainId+"/"+pathId
        }]
    });

    return false;
}

function delPathDetail(domainId, pathId){
    $.app.confirm(null, "确认删除域名路径映射"+pathId, function (){
        $.app.deleteJson(V3_API_URL + '/path/'+domainId+'/'+pathId, null, function (data) {
            $.app.show(data.msg);
            refreshPathList(domainId)
        })
    })

    return false;
}

function initFunction(){

    window.addrow = $.iGrid.EdatagridHandle.addRowHandle('#contactDg', function(){
    }, function(){
        return {"pointer_type":"Node","weight":1};
    }, true);


    window.deleterow = $.iGrid.EdatagridHandle.removeRowHandle('#contactDg', function(){
    }, function(rows, success){
        success.apply();
    }, true);

    window.save = $.iGrid.EdatagridHandle.saveHandle('#contactDg', function(){
    }, function(cs, success){
        success.apply();
    }, true);

    window.cancel = $.iGrid.EdatagridHandle.rollbackHandle('#contactDg', function(){
    }, function(cs, success){
        success.apply();
    }, true, true);
}

function openPathPatternDiag(){
    $.iDialog.openDialog({
        title: '测试匹配路径',
        minimizable:false,
        width: 680,
        height: 360,
        content: `
            <div style="margin: 0px;">
            </div>
            <div class="cubeui-fluid" id="company-detail-form">
                <div class="cubeui-row span5">

                    <div class="cubeui-col-sm12 ">
                        <label class="cubeui-form-label">匹配路径:</label>
                        <div class="cubeui-input-block">
                            <div style='width:100%;margin-bottom: 0px;margin-top: 0px;'>
                                <div class="cubeui-col-sm10" style='padding-right:1px'>
                                    <input type="text" data-toggle="cubeui-textbox" id="url_pattern" name="url_pattern"
                                       value='{{:url_pattern}}' data-options="required:false,prompt:'请求路径格式, 请求路径的进行PathMatcher匹配'">
                                </div>
                                <div class="cubeui-col-sm2" style='padding-right:1px'>
                                    <label class="cubeui-label">
                                    <a href="javascript:void(0)" data-toggle="cubeui-menubutton" data-options="{
                                            onClick:function(){
                                                $('#req_path').textbox('setValue', $('#url_pattern').textbox('getValue'))
                                                $.app.show('设置成功')
                                            },
                                            plain:true,
                                            iconCls:'fa fa-check-square-o',
                                            btnCls: 'cubeui-btn-orange'
                                        }">使用</a></label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="cubeui-row span5">

                    <div class="cubeui-col-sm12 ">
                        <label class="cubeui-form-label">测试路径:</label>
                        <div class="cubeui-input-block">
                            <div style='width:100%;margin-bottom: 0px;margin-top: 0px;'>
                                <div class="cubeui-col-sm10" style='padding-right:1px'>
                                    <input type="text" data-toggle="cubeui-textbox" id="test_url" name="test_url"
                                           value='{{:test_url}}' data-options="required:false,prompt:'测试请求路径,点击测试查看是否匹配'">
                                </div>
                                <div class="cubeui-col-sm2" style='padding-right:1px'>
                                    <label class="cubeui-label">
                                    <a href="javascript:void(0)" data-toggle="cubeui-menubutton" data-options="{
                                            onClick:function(){
                                                let p = {};
                                                p.test = $('#test_url').textbox('getValue')
                                                p.pattern = $('#url_pattern').textbox('getValue')
                                                $.app.post(V3_API_URL + '/path/pathpattern', p, function (data) {
                                                    $.app.info(data.msg)
                                                });
                                            },
                                            plain:true,
                                            iconCls:'fa fa-calculator',
                                            btnCls: 'cubeui-btn-blue'
                                        }">测试</a></label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            `,
        render:function(opts, handler){
            let v = $($.iTextbox.findInByName($("#pathDetailDlg"), "req_path")).textbox("getValue")
            handler.render({url_pattern:v})
        },
        buttonsGroup1: [{
            text: '使用',
            iconCls: 'fa fa-check-square-o',
            btnCls: 'cubeui-btn-blue',
            onClick: function (){
                alert(1)
                $.iDialog.closeOutterDialog($(this))
            }
        }]
    });
}

function openMatchPatternDiag(){
    $.iDialog.openDialog({
        title: '测试路径格式转换规则',
        minimizable:false,
        width: 680,
        height: 360,
        content: `
            <div style="margin: 0px;">
            </div>
            <div class="cubeui-fluid" id="company-detail-form">
                <div class="cubeui-row span5">

                    <div class="cubeui-col-sm12 ">
                        <label class="cubeui-form-label">路径格式:</label>
                        <div class="cubeui-input-block">
                            <div style='width:100%;margin-bottom: 0px;margin-top: 0px;'>
                                <div class="cubeui-col-sm10" style='padding-right:1px'>
                                    <input type="text" data-toggle="cubeui-textbox" id="match_pattern" name="match_pattern"
                                       value='{{:match_pattern}}' data-options="required:false,prompt:'请求路径格式, 请求路径的进行正则转换模式'">
                                </div>
                                <div class="cubeui-col-sm2" style='padding-right:1px'>
                                    <label class="cubeui-label">
                                    <a href="javascript:void(0)" data-toggle="cubeui-menubutton" data-options="{
                                            onClick:function(){
                                                $('#search_path').textbox('setValue', $('#match_pattern').textbox('getValue'))
                                                $.app.show('设置成功')
                                            },
                                            plain:true,
                                            iconCls:'fa fa-check-square-o',
                                            btnCls: 'cubeui-btn-orange'
                                        }">使用</a></label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="cubeui-row span5">
                    <div class="cubeui-col-sm12 ">
                        <label class="cubeui-form-label">转换格式:</label>
                        <div class="cubeui-input-block">
                            <div style='width:100%;margin-bottom: 0px;margin-top: 0px;'>
                                <div class="cubeui-col-sm10" style='padding-right:1px'>
                                    <input type="text" data-toggle="cubeui-textbox" id="match_format" name="match_format"
                                       value='{{:match_format}}' data-options="required:false,prompt:'转换格式，转换后的模式'">
                                </div>
                                <div class="cubeui-col-sm2" style='padding-right:1px'>
                                    <label class="cubeui-label">
                                    <a href="javascript:void(0)" data-toggle="cubeui-menubutton" data-options="{
                                            onClick:function(){
                                                $('#replace_path').textbox('setValue', $('#match_format').textbox('getValue'))
                                                $.app.show('设置成功')
                                            },
                                            plain:true,
                                            iconCls:'fa fa-check-square-o',
                                            btnCls: 'cubeui-btn-orange'
                                        }">使用</a></label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="cubeui-row span5">

                    <div class="cubeui-col-sm12 ">
                        <label class="cubeui-form-label">测试转换路径:</label>
                        <div class="cubeui-input-block">
                            <div style='width:100%;margin-bottom: 0px;margin-top: 0px;'>
                                <div class="cubeui-col-sm10" style='padding-right:1px'>
                                    <input type="text" data-toggle="cubeui-textbox" id="match_url" name="match_url"
                                           value='{{:match_url}}' data-options="required:false,prompt:'测试转换路径,点击测试查看转换路径结果'">
                                </div>
                                <div class="cubeui-col-sm2" style='padding-right:1px'>
                                    <label class="cubeui-label">
                                    <a href="javascript:void(0)" data-toggle="cubeui-menubutton" data-options="{
                                            onClick:function(){
                                                let p = {};
                                                p.test = $('#match_url').textbox('getValue')
                                                p.pattern = $('#match_pattern').textbox('getValue')
                                                p.format = $('#match_format').textbox('getValue')
                                                $.app.post(V3_API_URL + '/path/matchpattern', p, function (data) {
                                                    $.app.info('转换结果：'+data.data)                                                    
                                                    $('#match_result').textbox('setValue', data.data)
                                                });
                                            },
                                            plain:true,
                                            iconCls:'fa fa-calculator',
                                            btnCls: 'cubeui-btn-blue'
                                        }">测试</a></label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
                <div class="cubeui-row span5">
                    <div class="cubeui-col-sm12 ">
                        <label class="cubeui-form-label">转换路径结果:</label>
                        <div class="cubeui-input-block">
                            <div style='width:100%;margin-bottom: 0px;margin-top: 0px;'>
                                <div class="cubeui-col-sm10" style='padding-right:1px'>
                                    <input type="text" data-toggle="cubeui-textbox" id="match_result" name="match_result"
                                           value='' data-options="readonly:true,required:false,prompt:'转换路径结果,点击测试查看转换路径结果'">
                                </div>
                                <div class="cubeui-col-sm2" style='padding-right:1px'>
                                    <label class="cubeui-label">
                                    <a href="javascript:void(0)" data-toggle="cubeui-menubutton" data-options="{
                                            onClick:function(){
                                                $('#match_result').textbox('setValue', '')
                                            },
                                            plain:true,
                                            iconCls:'fa fa-trash-o',
                                            btnCls: 'cubeui-btn-red'
                                        }">测试</a></label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            `,
        render:function(opts, handler){
            let v = $($.iTextbox.findInByName($("#pathDetailDlg"), "search_path")).textbox("getValue")
            let v2 = $($.iTextbox.findInByName($("#pathDetailDlg"), "replace_path")).textbox("getValue")

            handler.render({match_pattern:v,match_format:v2})
        },
        buttonsGroup1: [{
            text: '使用',
            iconCls: 'fa fa-check-square-o',
            btnCls: 'cubeui-btn-blue',
            onClick: function (){
                alert(1)
                $.iDialog.closeOutterDialog($(this))
            }
        }]
    });
}