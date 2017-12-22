$(document).ready(function(){
	setButtonCSS();
});

//Click from Add Row Button, used to add a new row
function addFieldTypesTableRow(elem){
	$(elem).closest("li").append('<div style="margin-left:150px;"><textarea name="product_name[field][]"></textarea> <textarea name="product_name[value][]"></textarea>(<a onclick="deleteFieldTypesTableRow(this)" class="button">Delete</a>)</div>');
	setButtonCSS();
}

//Delete Row
function deleteFieldTypesTableRow(elem){
	$(elem).parent().remove();
}

//Update the Link CSS
function setButtonCSS(){
	$("a.button").css({
		"cursor":"pointer",
		"text-decoration":"underline"
		});
}
