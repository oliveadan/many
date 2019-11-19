function v(a) {
    if (window.confirm('投票之后不可更改，您确定要投票吗？')) {
        $.ajax({
            url: '/vote',
            type: "POST",
            data: {
                vote: a
            },
            success: function (info) {
                console.log(info);
                if (info.code == 1) {
                    $('#btn' + a).next().css("height", $('#btn' + a).next().height() + 1 + "px");

                    $('#vv1').css("margin-top", 300 - 20 - info.data.votes[0] +
                        "px");
                    $('#vv2').css("margin-top", 300 - 20 - info.data.votes[1] +
                        "px");
                    $('#vv3').css("margin-top", 300 - 20 - info.data.votes[2] +
                        "px");
                    $('#vv4').css("margin-top", 300 - 20 - info.data.votes[3] +
                        "px");
                    $('#vv5').css("margin-top", 300 - 20 - info.data.votes[4] +
                        "px");
                    $('#vv6').css("margin-top", 300 - 20 - info.data.votes[5] +
                        "px");
                    $('#vv7').css("margin-top", 300 - 20 - info.data.votes[6] +
                        "px");
                    $('#vv8').css("margin-top", 300 - 20 - info.data.votes[7] +
                        "px");
                    $('#vv9').css("margin-top", 300 - 20 - info.data.votes[8] +
                        "px");
                    $('#vv10').css("margin-top", 300 - 20 - info.data.votes[9] +
                        "px");
                    $('#vv11').css("margin-top", 300 - 20 - info.data.votes[10] +
                        "px");
                    $('#vv12').css("margin-top", 300 - 20 - info.data.votes[11] +
                        "px");

                    $('#btn' + a).next().find(".VoteSpan").html($('#btn' + a).next().height());
                    alert(info.msg);
                    setTimeout(function () {
                        $(".VoteValue").css("display", 'block');
                    }, 250);
                    setTimeout(function () {
                        $(".VoteSpan").css("display", 'block');
                    }, 500);

                    var VoteSpanTri = $("<span></span>");
                    VoteSpanTri.attr("class", "VoteSpanTri");
                    $('#btn' + a).next().find(".VoteSpan").append(VoteSpanTri);
                    $(".VoteText").attr("disabled", true);
                    $('#v1').html(info.data.votes[0] * 7);
                    $('#v2').html(info.data.votes[1] * 7);
                    $('#v3').html(info.data.votes[2] * 7);
                    $('#v4').html(info.data.votes[3] * 7);
                    $('#v5').html(info.data.votes[4] * 7);
                    $('#v6').html(info.data.votes[5] * 7);
                    $('#v7').html(info.data.votes[6] * 7);
                    $('#v8').html(info.data.votes[7] * 7);
                    $('#v9').html(info.data.votes[8] * 7);
                    $('#v10').html(info.data.votes[9] * 7);
                    $('#v11').html(info.data.votes[10] * 7);
                    $('#v12').html(info.data.votes[11] * 7);

                    $('#vv1').css('height', info.data.votes[0] + 'px');
                    $('#vv2').css('height', info.data.votes[1] + 'px');
                    $('#vv3').css('height', info.data.votes[2] + 'px');
                    $('#vv4').css('height', info.data.votes[3] + 'px');
                    $('#vv5').css('height', info.data.votes[4] + 'px');
                    $('#vv6').css('height', info.data.votes[5] + 'px');
                    $('#vv7').css('height', info.data.votes[6] + 'px');
                    $('#vv8').css('height', info.data.votes[7] + 'px');
                    $('#vv9').css('height', info.data.votes[8] + 'px');
                    $('#vv10').css('height', info.data.votes[9] + 'px');
                    $('#vv11').css('height', info.data.votes[10] + 'px');
                    $('#vv12').css('height', info.data.votes[11] + 'px');
                } else {
                    alert(info.msg);
                }
                var vul1 = $('#vv1').height();
                var vul2 = $('#vv2').height();
                var vul3 = $('#vv3').height();
                var vul4 = $('#vv4').height();
                var vul5 = $('#vv5').height();
                var vul6 = $('#vv6').height();
                var vul7 = $('#vv7').height();
                var vul8 = $('#vv8').height();
                var vul9 = $('#vv9').height();
                var vul10 = $('#vv10').height();
                var vul11 = $('#vv11').height();
                var vul12 = $('#vv12').height();
                if (vul1 > 0 && vul1 <= 3) {
                    $('#vv1').css('height', 3 + 'px');
                };
                if (vul2 > 0 && vul2 <= 3) {
                    $('#vv2').css('height', 3 + 'px');
                };
                if (vul3 > 0 && vul3 <= 3) {
                    $('#vv3').css('height', 3 + 'px');
                };
                if (vul4 > 0 && vul4 <= 3) {
                    $('#vv4').css('height', 3 + 'px');
                };
                if (vul5 > 0 && vul5 <= 3) {
                    $('#vv5').css('height', 3 + 'px');
                };
                if (vul6 > 0 && vul6 <= 3) {
                    $('#vv6').css('height', 3 + 'px');
                };
                if (vul7 > 0 && vul7 <= 3) {
                    $('#vv7').css('height', 3 + 'px');
                };
                if (vul8 > 0 && vul8 <= 3) {
                    $('#vv8').css('height', 3 + 'px');
                };
                if (vul9 > 0 && vul9 <= 3) {
                    $('#vv9').css('height', 3 + 'px');
                };
                if (vul10 > 0 && vul10 <= 3) {
                    $('#vv10').css('height', 3 + 'px');
                };
                if (vul11 > 0 && vul11 <= 3) {
                    $('#vv11').css('height', 3 + 'px');
                };
                if (vul12 > 0 && vul12 <= 3) {
                    $('#vv12').css('height', 3 + 'px');
                };
            }
        })
    }
};