$(function () {
    var $form = $('#form');
    $form.on('submit', function (e) {
        e.preventDefault();
    });
    $form.on('submit', function () {
        ret = request(
            'POST',
            $(this).attr('action'),
            $(this).serialize(),
            function (ret) {
                $.gritter.add({
                    title: 'Create Success',
                    time: 3000,
                    after_close: function () {
                        location.href = "/admin";
                    }
                });
            },
            function (ret) {
                $.gritter.add({
                    title: 'Create Failed',
                    text: ret.message || 'Unknown Error, Retry Later.',
                    time: 3000
                });
            }
        )
        return false;
    });
});

