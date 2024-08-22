document.getElementById('addSectionButton').addEventListener('click', function() {
    const textarea = document.getElementById('requestData');
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const selectedText = textarea.value.substring(start, end);
    const beforeText = textarea.value.substring(0, start);
    const afterText = textarea.value.substring(end);
    textarea.value = beforeText + '§' + selectedText + '§' + afterText;
});

document.getElementById('clearButton').addEventListener('click', function() {
const textarea = document.getElementById('requestData');
textarea.value = textarea.value.replace(/§/g, '');
});