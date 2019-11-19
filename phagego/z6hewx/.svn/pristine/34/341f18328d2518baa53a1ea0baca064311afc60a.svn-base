
var time=null;
var time2=null;


$('.cai').mouseenter(function(){
		$(".cai").animate({right:'0px'},"slow");
		$('.caini_xh').css('background','url(images/cai_04.png) no-repeat')
		clearTimeout(time)
		clearTimeout(time2)

	});

$('.cai').mouseleave(
	function(){
	time=setTimeout(function (){ 
			$(".cai").animate({right:'-460px'},"slow");
			$('.caini_xh').css('background','url(images/cai_03.png) no-repeat') 
},3000)

});

$(function(){
	setTimeout(function (){
		$(".cai").animate({right:'0px'},"slow");
		$('.caini_xh').css('background','url(images/cai_04.png) no-repeat')	
		time2=setTimeout(function (){ 
			$(".cai").animate({right:'-460px'},"slow");
			$('.caini_xh').css('background','url(images/cai_03.png) no-repeat') 
},5000)	

	},3000)


})





$('.caini_xh').click( function(){
		var ri=$('.cai').css('right')
		if(ri=='0px'){
			$(".cai").animate({right:'-460px'},"slow");
			$('.caini_xh').css('background','url(images/cai_03.png) no-repeat') 
		}
})





$(document).scroll(function() {
	var obody=$(window).scrollTop();
	if(obody>50){
		$(".xuanfu").show();
	}
	else{
		$(".xuanfu").hide();
	}
	
});

$(function(){
	var xian=$(window).height()-32;
	var wi=$(window).width();
	if(wi<1000){
	$('.navbar-nav').height(xian)
	}else{

	$('.navbar-nav').height(46)

	}
	var hei=$('.head').outerHeight()
	$('.zhan_wei').height(hei)
})









$('.zhanchu').click(function(){
	$(this).parent().toggleClass('chu')
})



$('.erweiwe').click(function(){

	$('.xufu_erw').toggleClass('show')

})