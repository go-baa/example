$(function () {
    var cur_page = 1;
    var more = 0;
    $post_list = $("#post-list");
    $prev = $("#pager").find(".previous");
    $next = $("#pager").find(".next");
    getlist(cur_page);

    function getlist(page) {
        var url = "/p/" + page;
        request(
            'GET',
            url,
            {},
            function (ret) {
                if (ret.data.items.length == 0) {
                    return false;
                }
                cur_page = ret.data.cur_page
                more = ret.data.more

                $post_list.empty()
                ret.data.items.forEach(function (item) {
                    item_temp = `
                        <div class="post-preview">
                            <a href="/c/${item.id}">
                                <h2 class="post-title">
                                    ${item.title}
                                </h2>
                            </a>
                            <p class="post-meta">Posted on ${item.update_time}</p>
                        </div>
                        <hr>
                    `;
                    itemDom = $(item_temp);
                    $post_list.append(itemDom);
                });
            },
            function (ret) {
                console.log(ret);
                return false;
            }
        )
        return false;
    }

    $next.find("a").on("click", function () {
        if (more == 0) {
            return false;
        }
        getlist(cur_page + 1);
        $('html,body').animate({ scrollTop: 0 }, 700)
    });

    $prev.find("a").on("click", function () {
        if (cur_page == 1) {
            return false;
        }
        getlist(cur_page - 1)
        $('html,body').animate({ scrollTop: 0 }, 700)
    });
});