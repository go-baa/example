$(function () {
    var $form = $('#login-form');
    $form.on('submit', function (e) {
        e.preventDefault();
    });
    $form.on('submit', function () {
        ret = request(
            'POST',
            $(this).attr('action'),
            $(this).serialize(),
            function (ret) {
                location.href = "/admin";
            },
            function (ret) {
                $.gritter.add({
                    title: 'Login Failed',
                    text: ret.message || 'Unknown Error, Retry Later.',
                    time: 3000
                });
            }
        )
        return false;
    });
});

