
$(function(){
	
	jQuery('.main-content').css({'min-height':$(window).height()});
	$.zpost('/',{},function(result){
		
		if (result&&result.ret==200&&!!result.data){
			var menu = '';
			$.each(result.data,function(i,m){
			 	if (m.ChildMenu && m.ChildMenu.length > 0) {
			 		menu+='<li class="menu-list"><a href="#" ><span>'+ m.DisplayName+'</span></a><ul class="sub-menu-list">';
			 		var subMenu='';
			 		$.each(m.ChildMenu,function(j,n){
			 			subMenu+='<li><a class="norefresh" href="'+n.ControlUrl+'"> '+n.DisplayName+'</a></li>';
			 		});
			 		menu+=subMenu+'</ul></li>';
				}else{
					menu+= '<li><a class="norefresh" href="'+m.ControlUrl+'" ><span>'+ m.DisplayName+'</span></a></li>';
				}
			});
			$('.js-left-nav').append(menu);
			leftSelect();
		}
	});

	//左边菜单加选中状态
	function leftSelect(){
		var pathname = location.pathname;
        if (pathname!="/"){
        	$('.js-left-nav .norefresh').filter(function(){           
	        	return $(this).attr('href') == pathname;
	        }).closest("li").addClass('active').parents('.menu-list').addClass('nav-active');
        }
	}

    $("body").delegate(".menu-list > a","click",function(){
        var parent = jQuery(this).parent();
        var sub = parent.find('> ul');
        //if(!jQuery('body').hasClass('left-side-collapsed')) {
         	if(sub.is(':visible')) {
         		parent.removeClass('nav-active');
	            // sub.slideUp(200, function(){
	            //    parent.removeClass('nav-active');
	            //    jQuery('.main-content').css({height: ''});
	            //    mainContentHeightAdjust();
	            // });
         	} else {
	            visibleSubMenuClose();
	            parent.addClass('nav-active');
	            sub.slideDown(200, function(){
	                //mainContentHeightAdjust();
            	});
        	}
        //}
      	return false;
    });

    var firstEnter = false;
	$("body").delegate(".js-left-nav .norefresh","click",function(){
		// if ($(this).closest("li").hasClass('active')){
		// 	return false;
		// }
		$('.js-left-nav .active').removeClass('active');
		$(this).closest("li").addClass('active');
		if(!jQuery('body').hasClass('left-side-collapsed')&&$('.js-left-nav .nav-active').length&&$(this).closest(".nav-active").length==0) {
			visibleSubMenuClose();
			// mainContentHeightAdjust();
		}
		var url=$(this).attr('href');
		getpage(url,firstEnter);
		firstEnter = true;
		// $('.wrapper').zload(url);
      	return false;
    });


	window.onpopstate = function (e) {
		if (e.state){
			$('.js-left-nav .active').removeClass('active');
			 visibleSubMenuClose();
			leftSelect();
			// mainContentHeightAdjust();
			$('.wrapper').empty().append(e.state.html);
			execjs(e.state.html);
		}
	}



    function visibleSubMenuClose() {
      	jQuery('.menu-list').each(function() {
         	var t = jQuery(this);
         	if(t.hasClass('nav-active')) {
	            t.find('> ul').slideUp(200, function(){
	               t.removeClass('nav-active');
	            });
         	}
      	});
    }

   	//function mainContentHeightAdjust() {
      // Adjust main content height
      	// var docHeight = jQuery(document).height();
      	// if(docHeight > jQuery('.main-content').height()){
       //   	jQuery('.main-content').height(docHeight);
      	// }
    //}
})

function getpage(url,isFirst){
	$.zget(url,"",function(result){
		history.pushState({html:result}, "what", url);
        $('.wrapper').empty().append(result);
        if (!isFirst) {
			execjs(result);
        }
	});
	return false;
}

$('html').on('click','.skip',function () {
	var href = $(this).attr('href');
    getpage(href);
    return false;
})

$('html').on('click','.skip1',function () {
	
	window.sessionStorage.URl = window.location.href;

    var imgType = $('#imgType').val();   //图片管理位置顺序
    window.sessionStorage.imgType = imgType;

    var agent_name = $('#agent_name').val();   //代理商名称
    window.sessionStorage.agent_name = agent_name;

	var jg_name = $('#jg_name').val();   //机构名称
	window.sessionStorage.jg_name = jg_name;

	var cooperationModel = $('#cooperationModel').val();   //合作类型
	window.sessionStorage.cooperationModel = cooperationModel;

	var loanType = $('#loanType').val();   //贷款类型
	window.sessionStorage.loanType = loanType;

    var productType = $('#productType').val();   //产品类型
    window.sessionStorage.productType = productType;


	var locationType = $('select[name="locationType"]').val();
	window.sessionStorage.locationType=locationType;

	var orderState = $('select[name="orderState"]').find('option:selected').text();   //订单状态
	window.sessionStorage.orderState=orderState;

	var product_name = $('#product_name').val();   //产品名称
	if(product_name != null){
		window.sessionStorage.product_name=product_name;
	}
	var agent_name = $('#agent_name').val();   //运营商名称
	if(agent_name != null){
		window.sessionStorage.agent_name=agent_name;
	}
	

	var agentName = $('#agentName').val();   //代理商名称
	window.sessionStorage.agentName=agentName;

	var push_text = $('#push_text').val();   //推送文本
	window.sessionStorage.push_text=push_text;

	var send_text = $('#send_text').val();   //发送文本
	window.sessionStorage.send_text=send_text;
	
	
	var proName = $('#proName').val();   //产品名称
	window.sessionStorage.proName=proName;

	var user_name = $('.user_name').val();   //用户名
	window.sessionStorage.user_name=user_name;
	
	var phone_number = $('.phone_number').val();   //手机号码
	window.sessionStorage.phone_number = phone_number;

	var id_number = $('.id_number').val();   //用户身份证号码
	window.sessionStorage.id_number=id_number;

	var start_time = $('#startDate').val();   //开始时间
	window.sessionStorage.start_time = start_time;

	var end_time  =$('#endDate').val();      //结束时间
	window.sessionStorage.end_time = end_time;
	
    //已结算平台
    var sort = $('#sort').val();   
    window.sessionStorage.sort = sort;

    //未结算平台
    var start_time1 = $('#startDate1').val();   //开始时间
    window.sessionStorage.start_time1 = start_time1;

    var sort1 = $('#sort1').val();   //时间排序
    window.sessionStorage.sort1 = sort1;
    
    var href = $(this).attr('href');
	getpage(href);
	return false;
});

$('html').on('click','.skip2',function () {
	window.sessionStorage.URl = window.location.href;
	var href = $(this).attr('href');
	getpage(href);
	return false;
});

$('html').on('click','.skip3',function () {

	window.sessionStorage.URl = window.location.href;

    var agent_name = $('#agent_name').val();   //代理商名称
    window.sessionStorage.agent_name = agent_name;

	return false;
});

$('html').on('mousedown','.norefresh',function () {
    window.sessionStorage.clear();

})
$('html').on('mouseup','.pull-left button[type="submit"]',function () {
    window.sessionStorage.clear();
})
// function zaq(name) {
//     var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
//     var r = window.location.search.substr(1).match(reg);
//     if(r!=null)return  unescape(r[2]); return null;
// }


function execjs(html){
	// 第一步：匹配加载的页面中是否含有js
	var regDetectJs = /<script(.|\n)*?>(.|\n|\r\n)*?<\/script>/ig;
	var jsContained = html.match(regDetectJs);
	// 第二步：如果包含js，则一段一段的取出js再加载执行
	if(jsContained) {
		// 分段取出js正则
		var regGetJS = /<script(.|\n)*?>((.|\n|\r\n)*)?<\/script>/im;

		// 按顺序分段执行js
		var jsNums = jsContained.length;
		for (var i=0; i<jsNums; i++) {
			var jsSection = jsContained[i].match(regGetJS);
			if(jsSection[2]) {
				if(window.execScript) {
					// 给IE的特殊待遇
					window.execScript(jsSection[2]);
				} else {
					// 给其他大部分浏览器用的
					window.eval(jsSection[2]);
				}
			}
		}
	}
}
