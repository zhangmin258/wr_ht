jQuery.zget = function(url, args, callback) {
	if (args){
		if (typeof(args)=='string'){
			args+="&pushstate=1"
		}else if (typeof(args)=='object'){
			args.pushstate=1
		}
	}else{
		args={pushstate:1}
	}
	
	jQuery.get(url, args, function(result){
		if (result.ret==408){
			location.href='/login'
		}else{
			callback(result)
		}
	})
}

jQuery.fn.extend.zload = function(url) {
	if (url.indexOf("?")==-1){
		url+="?pushstate=1"
	}else if (url.lastIndexOf("?")==url.length-1){
		url+="pushstate=1"
	}
	else{
		url+="&pushstate=1"
	}
	jQuery(this).load(url)
}

jQuery.zpost = function(url, args, callback) {
	// args.pushstate=1
	jQuery.post(url, args, function(result){
		if (result.ret==408){
			location.href='/login'
		}else{
			callback(result)
		}
	})
}