$(function() {
   var editor = new Simditor({
        textarea: $('#editor'),
        //optional options
        placeholder: 'Input content',
        upload: {
            url: '/admin/upload',
            params: null,
            fileKey: 'upload_file',
            connectionCount: 3,
            leaveConfirm: 'Uploading is in progress, are you sure to leave this page?'
        },
        tabIndent: true,
        pasteImage: true,
        cleanPaste: false,
        toolbarFloat: true,
        toolbarFloatOffset: 50,
        // imageButton: ['upload'],
        codeLanguages: [
            { name: 'Bash', value: 'bash' },
            { name: 'C++', value: 'c++' },
            { name: 'C#', value: 'cs' },
            { name: 'CSS', value: 'css' },
            { name: 'Erlang', value: 'erlang' },
            { name: 'Less', value: 'less' },
            { name: 'Sass', value: 'sass' },
            { name: 'Diff', value: 'diff' },
            { name: 'CoffeeScript', value: 'coffeescript' },
            { name: 'HTML,XML', value: 'html' },
            { name: 'JSON', value: 'json' },
            { name: 'Java', value: 'java' },
            { name: 'JavaScript', value: 'js' },
            { name: 'Markdown', value: 'markdown' },
            { name: 'Objective C', value: 'oc' },
            { name: 'PHP', value: 'php' },
            { name: 'Perl', value: 'parl' },
            { name: 'Python', value: 'python' },
            { name: 'Ruby', value: 'ruby' },
            { name: 'SQL', value: 'sql'}
        ]
    });
});