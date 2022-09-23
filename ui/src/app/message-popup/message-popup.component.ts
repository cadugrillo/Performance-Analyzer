import { Component, Inject } from '@angular/core';
import {MAT_DIALOG_DATA} from '@angular/material/dialog';

export interface DialogData {
  title: string;
  text: string;
}

@Component({
  selector: 'app-message-popup',
  templateUrl: './message-popup.component.html',
  styleUrls: ['./message-popup.component.css']
})

export class MessagePopupComponent {

  constructor(@Inject(MAT_DIALOG_DATA) public data: DialogData) {}
}
