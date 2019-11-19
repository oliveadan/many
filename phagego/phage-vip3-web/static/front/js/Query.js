// 弹出
function showDetail() {

    $.ajax({
        url: '/query',
        dataType: 'json',
        cache: false,
        data: {
            account: $("#search1").val()
        },
        type: 'Post',
        success: function (data) {
            if (data.code === 2) {
                layer.msg("会员账号不存在");
                $('body #show').hide();
                return
            }
            var level = "";
            var lucky = "";
            //总信息
            var total = "";
            total = "<div style=\"width:13%\">" + data.Account + "</div>"
                + "<div>" + data.level + "</div>"
                + "<div style=\"width:20%\">" + data.membertotalbet + "</div>"
                + "<div style=\"width:20%\">" + data.totallevelgift + "</div>"
                + "<div style=\"width:10%\">" + data.totalluckygift + "</div>"
                + "<div style=\"width:25%\">" + data.balance + "</div>";
            //周信息
            var week = "";
            $.each(data.membersingle, function (i, award) {
                var account = award.Account;
                var id = award.Id
                var levelgift = award.LevelGift;
                var luckygift = award.LuckyGift;
                var levelid = account + levelgift;
                var luckyid = levelgift + "s" + luckygift;
                if (award.LevelGift === 0) {
                    level = "<div style=\"width:9%\"><button disabled>无奖励</button></div>"
                } else if (award.LevelEnable === 1) {
                    level = "<div style=\"width:9%\"><button disabled>已领取</button></div>"
                } else {
                    level = '<div id=' + levelid + ' style="width:9%">' + '<button  onclick="GetLevelGift(' + id + ',' + levelgift + ',\'' + account + '\')">' + '点击领取' + '</button>' + '</div>'
                }
                if (award.LuckyGift === 0) {
                    lucky = "<div  style=\"width:10%\"><button disabled>无奖励</button></div>"
                } else if (award.LuckyEnable === 3) {
                    lucky = "<div  style=\"width:10%\"><button disabled>已过期</button></div>"
                } else if (award.LuckyEnable === 1) {
                    lucky = "<div  style=\"width:10%\"><button disabled>已领取</button></div>"
                } else {
                    lucky = '<div id=' + luckyid + ' style="width:10%">' + '<button  onclick="GetLuckyGift(' + levelgift + ',' + id + ',' + luckygift + ',\'' + account + '\')">' + '点击领取' + '</button>' + '</div>'
                }
                x = "<div style=\"width:13%\">" + award.Account + "</div>"
                    + "<div style=\"width:15%\">" + award.Bet + "</div>"
                    + "<div style=\"width:18%\">" + award.LevelGift + "</div>"
                    + level
                    + "<div style=\"width:15%\">" + award.LuckyGift + "</div>"
                    + lucky
                    + "<div style=\"width:20%\">" + award.PeriodName + "</div>";
                week += x;
            });
            if (data.Account != "") {
                $("#membertotal").html(total);
            }
            if (data.membersingle.length != 0) {
                $("#weekup").html(week);
                str = "共" + data.membersingle.length + "条";
                $("#countweek").text(str)
            }
        }
    });
    $('body #show').show();
};
$('#close').click(function () {
    $(this).parent().parent().hide();
});

function GetLevelGift(id, i, h) {
    $.ajax({
        url: '/getgift',
        dataType: 'json',
        cache: false,
        type: 'POST',
        data: {id: id, account: h, gift: i, type: 0},
        beforeSend: function (){
            layer.load()
        },
        success: function (obj) {
            layer.closeAll('loading');
            if (obj.code === 1) {
                s = "#" + h + i;
                $(s).html('<button disabled>已领取</button>')
            }
            layer.msg(obj.msg);
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            var x = 1;
            localStorage.clear();
        }
    });
}

function GetLuckyGift(l, id, i, h) {
    $.ajax({
        url: '/getgift',
        dataType: 'json',
        cache: false,
        type: 'POST',
        data: {id: id, account: h, gift: i, type: 1},
        beforeSend: function (){
            layer.load()
        },
        success: function (obj) {
            layer.closeAll('loading');
            if (obj.code === 1) {
                s = "#" + l + "s" + i;
                $(s).html('<button disabled>已领取</button>')
            }
            layer.msg(obj.msg);
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            var x = 1;
            localStorage.clear();
        }
    });
}