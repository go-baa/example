$(function () {
    $talbe = $("#data-table");
    var $form = $('form');
    var $filterKeyword = $form.find('[name=keyword]');
    var filterParams = {};
    $filterKeyword.on('change', search);

    function search() {
        table.draw();
    }

    $form.on('submit', function (e) {
        e.preventDefault();
        search();
    });

    table = $talbe.DataTable({
        ordering: false,
        searching: false,
        lengthChange: false,
        pageLength: 10,
        responsive: true,
        processing: true,
        serverSide: true,
        rowId: 'id',
        ajax: {
            url: "/admin/content/search",
            data: function (params) {
                delete params.columns;
                delete params.search;
                params["keyword"] = $filterKeyword.val()
            },
            dataSrc: function(ret) {
                if (!ret || ret.state) {
                    if (ret.message) {
                        ret.error = ret.message;
                    }
                    return [];
                }

                ret.data || (ret.data = {});
                ret.recordsTotal = ret.data.total || 0;
                ret.recordsFiltered = ret.data.total || 0;

                return ret.data.data || [];
            }
        },
        columns: [
            { "data": "id" },
            { "data": "title" },
            { "data": "update_time_full" }
        ],
        columnDefs: [
            {
                'render': function (data, type, row) {
                    return [
                        '<a class="m-r-5" target="_blank" href="/c/' + row.id + '">View</a>',
                        '<a class="m-r-5" href="/admin/content/' + row.id + '/edit">Edit</a>',
                        '<a class="m-r-5 action-delete" href="#">Delete</a>',
                    ].join('');
                },
                'targets': 3
            }
        ],
        createdRow: function (row, data, index) {
            var $row = $(row);
            $row.find('.action-delete').on('click', function () {
                delmodal = $("#delete-modal");
                delmodal.modal("show");
                delmodal.find("#delete-confirm").on("click", function() {
                    delmodal.modal("hide");
                    request(
                        'DELETE',
                        '/admin/content/' + data.id,
                        {},
                        function (ret) {
                            $.gritter.add({
                                title: 'Delete Success',
                                time: 3000,
                                after_close: function() {
                                    table.draw();
                                }
                            });
                        },
                        function (ret) {
                            $.gritter.add({
                                title: 'Delete Failed',
                                text: ret.message || 'Unknown Error, Retry Later.',
                                time: 3000
                            });
                        }
                    )
                    return false;
                })
            });
        }
    });
});