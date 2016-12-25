$(function () {
    var $form = $('#chgpwd-form');
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
                    title: 'Update Success',
                    time: 1000,
                    after_close: function() {
                        history.go(-1);
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

