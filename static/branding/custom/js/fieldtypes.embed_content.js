$(document).ready(function(){

	$( ".embed_content tbody" ).sortable();
});

//Click from Add Linkage Button, used to add a new row
function addEmbedContentLinkage(elem){

	$(elem).parent().find("table tbody").append('<tr><td><span class="ui-icon ui-icon-arrowthick-2-n-s"></span></td><td>'+$(elem).closest("li").find("textarea.hidden_content_option").val()+'</td><td><a target="_blank" class="button embed_url" href="">View</a></td><td><a onclick="deleteEmbedContent(this)" class="button">Delete</a></td></tr>');
}

//Delete Row
function deleteEmbedContent(elem){
	$(elem).closest("tr").remove();
}

//After user select the different content, refresh the preview link
function changeEmbedContentLink(elem){
	$(elem).closest("tr").find(".embed_url").attr("href", $(elem).find("option:selected").attr("url"));
}

//Redirect to add content page
function addEmbedContent(type_id){
	window.open("admincp/publish/create_post/"+type_id);
}