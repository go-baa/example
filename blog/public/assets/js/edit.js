$(function () {
    var $form = $('#form');
    $form.on('submit', function (e) {
        e.preventDefault();
    });
    cid = $('#cid').val();
    $form.on('submit', function () {
        ret = request(
            'PUT',
            '/admin/content/' + cid,
            $(this).serialize(),
            function (ret) {
                $.gritter.add({
                    title: 'Update Success',
                    time: 3000,
                    after_close: function () {
                        location.href = "/admin";
                    }
                });
            },
            function (ret) {
                $.gritter.add({
                    title: 'Update Failed',
                    text: ret.message || 'Unknown Error, Retry Later.',
                    time: 3000
                });
            }
        )
        return false;
    });
});

