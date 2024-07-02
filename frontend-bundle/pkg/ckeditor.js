/**
 * @license Copyright (c) 2014-2024, CKSource Holding sp. z o.o. All rights reserved.
 * For licensing, see LICENSE.md or https://ckeditor.com/legal/ckeditor-oss-license
 */

import { ClassicEditor } from '@ckeditor/ckeditor5-editor-classic';

import { Alignment } from '@ckeditor/ckeditor5-alignment';
import { Autoformat } from '@ckeditor/ckeditor5-autoformat';
import { Bold, Italic } from '@ckeditor/ckeditor5-basic-styles';
import { BlockQuote } from '@ckeditor/ckeditor5-block-quote';
import { CKBox } from '@ckeditor/ckeditor5-ckbox';
import { CloudServices } from '@ckeditor/ckeditor5-cloud-services';
import { Heading } from '@ckeditor/ckeditor5-heading';
import {
  Image,
  ImageCaption,
  ImageStyle,
  ImageToolbar,
  ImageResize,
  ImageUpload,
  PictureEditing,
  ImageInsert,
  ImageInsertViaUrl,
  ImageBlock,
} from '@ckeditor/ckeditor5-image';
import { Indent, IndentBlock } from '@ckeditor/ckeditor5-indent';
import { LinkImage } from '@ckeditor/ckeditor5-link';
import { List } from '@ckeditor/ckeditor5-list';
import { MediaEmbed } from '@ckeditor/ckeditor5-media-embed';
import { Paragraph } from '@ckeditor/ckeditor5-paragraph';
import { PasteFromOffice } from '@ckeditor/ckeditor5-paste-from-office';
import { Table, TableToolbar } from '@ckeditor/ckeditor5-table';
import { TextTransformation } from '@ckeditor/ckeditor5-typing';
import { Undo } from '@ckeditor/ckeditor5-undo';

// You can read more about extending the build with additional plugins in the "Installing plugins" guide.
// See https://ckeditor.com/docs/ckeditor5/latest/installation/plugins/installing-plugins.html for details.

class Editor extends ClassicEditor {
  static builtinPlugins = [
    Alignment,
    Heading,
    Autoformat,
    BlockQuote,
    Bold,
    Image,
    ImageCaption,
    ImageStyle,
    ImageToolbar,
    ImageResize,
    ImageUpload,
    ImageInsert,
    ImageBlock,
    ImageInsertViaUrl,
    CKBox,
    CloudServices,
    LinkImage,
    Indent,
    IndentBlock,
    Italic,
    List,
    MediaEmbed,
    Paragraph,
    PasteFromOffice,
    PictureEditing,
    Table,
    TableToolbar,
    TextTransformation,
    Undo,
  ];

  //@ts-ignore
  static defaultConfig = {
    toolbar: {
      items: [
        'heading',
        '|',
        'bold',
        'italic',
        'bulletedList',
        'numberedList',
        '|',
        'linkImage',
        '|',
        'outdent',
        'indent',
        '|',
        'alignment',
        '|',
        'blockQuote',
        'insertImage',
        'insertTable',
        'mediaEmbed',
        '|',
        'undo',
        'redo',
      ],
    },
    indentBlock: {
      offset: 1,
      unit: 'em',
    },
    image: {
      toolbar: [
        {
          title: 'Image Style',
          name: 'imageStyle:customDropdown',
          items: [
            'imageStyle:alignLeft',
            'imageStyle:alignRight',
            'imageStyle:alignCenter',
            'imageStyle:alignBlockRight',
            'imageStyle:alignBlockLeft',
          ],
          defaultItem: 'imageStyle:alignLeft',
        },
        '|',
        'toggleImageCaption',
        'imageTextAlternative',
        'ImageResize',
        '|',
        'linkImage',
      ],
      insert: {
        // If this setting is omitted, the editor defaults to 'block'.
        // See explanation below.
        integrations: ['url'],
        type: 'auto',
      },
    },
    table: {
      contentToolbar: ['tableColumn', 'tableRow', 'mergeTableCells'],
    },
  };
}

export default Editor;
