function request (method, url, data, success, error) {
    return $.ajax({
        type: method,
        url: url,
        data: data,
        success: function (ret) {
            if (ret && ret.code == 0) {
                success && success(ret);
            } else {
                error && error(ret);
            }
        },
        dataType: 'json'
    });
}