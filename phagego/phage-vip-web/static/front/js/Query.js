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
            //总信息
            var total = "";
            total = "<div>" + data.Account + "</div>"
                + "<div>" + data.Level + "</div>"
                + "<div style=\"width:20%\">" + data.membertotalbet + "</div>"
                + "<div style=\"width:20%\">" + data.totallevelgift + "</div>"
                + "<div style=\"width:10%\">" + data.totalluckygift + "</div>"
                + "<div style=\"width:25%\">" + data.balance + "</div>";
            //周信息
            var week = "";
            $.each(data.membersingle, function (i, award) {
                x = "<div style=\"width:20%\">" + award.Account + "</div>"
                    + "<div style=\"width:20%\">" + award.Bet + "</div>"
                    + "<div style=\"width:20%\">" + award.LevelGift + "</div>"
                    + "<div style=\"width:20%\">" + award.LuckyGift + "</div>"
                    + "<div style=\"width:20%\">" + award.PeriodName + "</div>";
                week += x;
            });
             if (data.Account != "") {
                $("#membertotal").html(total);
            }
            if (data.membersingle.length != 0) {
                $("#weekup").html(week)
                str = "共"+data.weekup.length+"条"
                $("#countweek").text(str)
            }
        }
    });
    $('body #show').show();
};
$('#close').click(function () {
    $(this).parent().parent().hide();
});