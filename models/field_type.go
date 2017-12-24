package models

func GetFieldTypeHTML(fieldType string) string {
	if fieldType == "date" {
		return `<li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="future_only">Future Only</label>
						<input type="checkbox" name="future_only" value="1" class="checkbox" />
						<div class="help">Only allow future dates?</div>
					</li>`
	} else if fieldType == "checkbox" {
		return `<li>
						<label for="default">Default State</label>
						<select name="default">
<option value="checked">Checked</option>
<option value="unchecked" selected="selected">Unchecked</option>
</select>

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this box must be checked for the form to be processed.</div>
					</li>`
	} else if fieldType == "datetime" {
		return `<li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="future_only">Future Only</label>
						<input type="checkbox" name="future_only" value="1" class="checkbox" />
						<div class="help">Only allow future dates?</div>
					</li>`
	} else if fieldType == "embed_content" {
		return `<li>
						<label for="content_type">Content Type</label>
						<select name="content_type">
<option value="11">Blog</option>
<option value="3">Events</option>
<option value="13">name2</option>
<option value="12">names</option>
<option value="10">News</option>
</select>

					</li><li>
						<label for="field_name">Field Name</label>
						<input type="text" name="field_name" id="field_name" value="title" style="width: 275px" class="required text" />
						<div class="help">Enter the "system_name" for the field that you want to use to identify the content items.  For example, for a content type
	               		   like news articles, you might enter "headline" if you have a Headline field.  You can confirm the accuracy of this fieldname
	               		   by going to Publish > Content Types > Manage Fields (for the content type you have selected above).</div>
					</li><li>
						<label for="allow_multiple">Allow Multiple Relationships</label>
						<input type="checkbox" name="allow_multiple" value="1" class="checkbox" />
						<div class="help">If checked, the user can select one or many content items from the list.</div>
					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li>`
	} else if fieldType == "file_list" {
		return `<li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="multi_files">Allow Multiple Files</label>
						<input type="checkbox" name="multi_files" value="1" class="checkbox" />
						<div class="help">Allow user select more than 1 files</div>
					</li>`
	} else if fieldType == "file_upload" {
		return `<li>
						<label for="filetypes">Allowed Filetypes</label>
						<input type="text" name="filetypes" id="filetypes" value="" style="width: 275px" class="text" />
						<div class="help">Enter the filetypes (e.g., "jpg", "gif", "pdf", and "doc") that can be uploaded here.  Though not a foolproof mechanism
	          	      for validating filetypes, validating the file extension will help make sure people upload proper files here.  If someone
	          	      does upload a malicious file by renaming the file, the file will still be non-executable as all filenames are encrypted and
	          	      securely stored.</div>
					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="width">Width</label>
						<input type="text" name="width" id="width" value="250px" style="width: 75px" class="text" />
						<div class="help">Enter the width of this field.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, a file must be uploaded here for the form to be processed.</div>
					</li>`
	} else if fieldType == "gallery" {
		return `<li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="multi_files">Allow Multiple Files</label>
						<input type="checkbox" name="multi_files" value="1" class="checkbox" />
						<div class="help">Allow user select more than 1 files</div>
					</li>`
	} else if fieldType == "member_group_relationship" {
		return `<li>
						<label for="allow_multiple">Allow Multiple Relationships</label>
						<input type="checkbox" name="allow_multiple" value="1" class="checkbox" />
						<div class="help">If checked, the user can select one or many member groups from the list.</div>
					</li><li>
						<label for="default">Default Value</label>
						<select name="default">
<option value="1">默认</option>
</select>

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, a selection must be made form to be processed.</div>
					</li>`
	} else if fieldType == "multicheckbox" {
		return `<li>
						<label for="options">Options</label>
						<textarea name="options" style="width: 500px; height: 150px" class="required textarea"></textarea>
						<div class="help">Enter each option on a newline. If you want the option database value to be different than the option the user actually selects, enter it in the format of "Name=Value".</div>
					</li><li>
						<label for="default">Default Selection(s)</label>
						<textarea name="default" style="width: 275px; height: 80px" class="textarea"></textarea>
						<div class="help">To select multiple values by default, enter each value on a newline.</div>
					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, at least one checkbox must be checked a successful form submission.</div>
					</li>`
	} else if fieldType == "multiselect" {
		return `<li>
						<label for="options">Options</label>
						<textarea name="options" style="width: 500px; height: 150px" class="required textarea"></textarea>
						<div class="help">Enter each option on a newline. If you want the option database value to be different than the option the user actually selects, enter it in the format of "Name=Value".</div>
					</li><li>
						<label for="default">Default Selection(s)</label>
						<textarea name="default" style="width: 275px; height: 80px" class="textarea"></textarea>
						<div class="help">To select multiple values by default, enter each value on a newline.</div>
					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this field must not be empty for a successful form submission.</div>
					</li>`
	} else if fieldType == "radio" {
		return `<li>
						<label for="options">Options</label>
						<textarea name="options" style="width: 500px; height: 150px" class="required textarea"></textarea>
						<div class="help">Enter each option on a newline. If you want the option database value to be different than the option the user actually selects, enter it in the format of "Name=Value".</div>
					</li><li>
						<label for="default">Default Selection</label>
						<input type="text" name="default" id="default" value="" style="width: 275px" class="text" />

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this field must not be empty for a successful form submission.</div>
					</li>`
	} else if fieldType == "select" {
		return `<li>
						<label for="options">Options</label>
						<textarea name="options" style="width: 500px; height: 150px" class="required textarea"></textarea>
						<div class="help">Enter each option on a newline. If you want the option database value to be different than the option the user actually selects, enter it in the format of "Name=Value".</div>
					</li><li>
						<label for="default">Default Selection</label>
						<input type="text" name="default" id="default" value="" style="width: 275px" class="text" />

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this field must not be empty for a successful form submission.</div>
					</li>`
	} else if fieldType == "table" {
		return `<li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li>`
	} else if fieldType == "text" {
		return `<li>
						<label for="default">Default Value</label>
						<input type="text" name="default" id="default" value="" style="width: 275px" class="text" />

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="width">Width</label>
						<input type="text" name="width" id="width" value="250px" style="width: 75px" class="text" />
						<div class="help">Enter the width of this field.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this box must not be empty for the form to be processed.</div>
					</li><li>
						<label for="validators">Validators</label>
						<div style="float: left"><div class="check_option"><input type="checkbox" name="validators[]" value="trim"  /> Trim whitespace from around response</div><div class="check_option"><input type="checkbox" name="validators[]" value="strip_tags"  /> Strip HTML tags</div><div class="check_option"><input type="checkbox" name="validators[]" value="alpha_numeric"  /> Only alphanumeric characters</div><div class="check_option"><input type="checkbox" name="validators[]" value="numeric"  /> Only numbers</div><div class="check_option"><input type="checkbox" name="validators[]" value="valid_domain"  /> Must be a valid domain (e.g., "yahoo.com")</div><div class="check_option"><input type="checkbox" name="validators[]" value="valid_email"  /> Must be a valid email address (e.g., "test@example.com")</div>				</div>

					</li>`
	} else if fieldType == "textarea" {
		return `<li>
						<label for="default">Default Value</label>
						<textarea name="default" style="width: 500px; height: 80px" class="textarea"></textarea>

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="width">Width</label>
						<input type="text" name="width" id="width" value="250px" style="width: 75px" class="text" />
						<div class="help">Enter the width of this field.</div>
					</li><li>
						<label for="height">Height</label>
						<input type="text" name="height" id="height" value="80px" style="width: 75px" class="text" />
						<div class="help">Enter the height of this field.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this box must not be empty for the form to be processed.</div>
					</li><li>
						<label for="validators">Validators</label>
						<div style="float: left"><div class="check_option"><input type="checkbox" name="validators[]" value="trim"  /> Trim whitespace from around response</div><div class="check_option"><input type="checkbox" name="validators[]" value="strip_tags"  /> Strip HTML tags</div><div class="check_option"><input type="checkbox" name="validators[]" value="alpha_numeric"  /> Only alphanumeric characters</div>				</div>

					</li>`
	} else if fieldType == "topic_relationship" {
		return `<li>
						<label for="allow_multiple">Allow Multiple Relationships</label>
						<input type="checkbox" name="allow_multiple" value="1" class="checkbox" />
						<div class="help">If checked, the user can select one or many topics from the list.</div>
					</li><li>
						<label for="default">Default Selection(s)</label>
						<select name="default[]"  multiple="multiple">
<option value="" selected="selected"></option>
<option value="1011">军事</option>
<option value="1012">政治</option>
<option value="1010">热点</option>
<option value="1000">默认话题</option>
</select>

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, a selection must be made form to be processed.</div>
					</li>`
	} else if fieldType == "wysiwyg" {
		return `<li>
						<label for="default">Default Value</label>
						<textarea name="default" style="width: 500px; height: 80px" class="textarea"></textarea>

					</li><li>
						<label for="help">Help Text</label>
						<textarea name="help" style="width: 500px; height: 80px" class="textarea"></textarea>
						<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
					</li><li>
						<label for="required">Required Field</label>
						<input type="checkbox" name="required" value="1" class="checkbox" />
						<div class="help">If checked, this box must not be empty for the form to be processed.</div>
					</li><li>
						<label for="use_basic">Use Basic Editor</label>
						<input type="checkbox" name="use_basic" value="1" class="checkbox" />
						<div class="help">The "Basic" editor doesn't have all of the features of the WYSIWYG editor, but is more appropriate when you just want
	   		          basic HTML stylings, images, links, etc.</div>
					</li><li>
						<label for="validators">Validators</label>
						<div style="float: left"><div class="check_option"><input type="checkbox" name="validators[]" value="trim"  /> Trim whitespace from around response</div>				</div>

					</li>`
	} else if fieldType == "" {
		return ``
	}

	return ``
}
