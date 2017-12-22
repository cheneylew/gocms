$(document).ready(function(){
	$( ".upload_gallery" ).sortable({
		placeholder: "ui-state-highlight"
	});
	$( ".file_list tbody" ).sortable();
});

//Open KCFinder Browser from image
function openKCFinderImage(n, multi_files) {
	if(multi_files != "1" && $("#"+n).find("li").length>0){
		alert("Only allow to select 1 image, please delete the selected image first");
		return ;
	}
    window.KCFinder = {};
    window.KCFinder.callBack = function(url) {
		addGalleyImage(url,n);
        // Actions with url parameter here
        window.KCFinder = null;
    };
    window.open('branding/default/js/ckeditor/kcfinder/browse.php?type=images&opener=custom&langCode=en', 'kcfinder_single');
}

//After insert the image from KCFinder
function addGalleyImage(path,n){
	$("#"+n).append("<li><input name='"+n+"[]' type='hidden' value='"+path+"'/><img src='"+path+"'/> (<a onclick='deleteUploadImage(this)' class='button'>Delete</a>)</li>");
}

//Open KCFinder Browser from file
function openKCFinderLink(n, multi_files) {
	if(multi_files != "1" && $("#"+n).find("tbody tr").length>0){
		alert("Only allow to select 1 file, please delete the selected file first");
		return ;
	}
	
    window.KCFinder = {};
    window.KCFinder.callBack = function(url) {
		addFileList(url,n);
        // Actions with url parameter here
        window.KCFinder = null;
    };
    window.open('branding/default/js/ckeditor/kcfinder/browse.php?type=files&opener=custom&langCode=en', 'kcfinder_single');
}

//After insert the file from KCFinder
function addFileList(path,n){
	
	$("#"+n).append("<tr><td><span class='ui-icon ui-icon-arrowthick-2-n-s'></span></td><td><input name='"+n+"[title][]' type='text' value='"+(path).replace(/^.*[\\\/]/, '')+"'/></td><td><input name='"+n+"[path][]' type='hidden' value='"+path+"'/><a href='"+path+"'>"+(path).replace(/^.*[\\\/]/, '')+"</a></td><td><a onclick='deleteUploadFile(this)' class='button'>Delete</a></td></tr>");
}

//Delete selected item image
function deleteUploadImage(elem){
	$(elem).closest("li").remove();
}

//Delete selected item file
function deleteUploadFile(elem){
	$(elem).closest("tr").remove();
}
