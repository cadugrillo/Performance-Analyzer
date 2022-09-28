import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { MessagePopupComponent } from '../message-popup/message-popup.component';
import { WaitPopupComponent } from '../wait-popup/wait-popup.component';

import { IUser, CognitoService } from '../../services/cognito.service';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-forgot-password',
  templateUrl: './forgot-password.component.html',
  styleUrls: ['./forgot-password.component.css']
})
export class ForgotPasswordComponent implements OnInit {

  loading: boolean;
  isConfirm: boolean;
  user: IUser;

  constructor(private router: Router, public dialog: MatDialog,
    private cognitoService: CognitoService) {
      this.loading = false;
      this.isConfirm = false;
      this.user = {} as IUser;
     }

  ngOnInit(): void {
  }

  public forgotPassword() {
    if (this.user.email != null) {
      this.dialog.open(WaitPopupComponent, {});
      this.cognitoService.fogotPassword(this.user).then(() => {
        this.dialog.closeAll();
        this.isConfirm = true;
      }).catch(() => {
        this.dialog.closeAll();
        this.dialog.open(MessagePopupComponent, {data: {title: "Error", text: "No user registered with this email!"}});
      })
    } else this.dialog.open(MessagePopupComponent, {data: {title: "Error", text: "Empty email!"}});
  }

  public confirmNewPassword() {
    this.dialog.open(WaitPopupComponent, {});
    this.cognitoService.confirmNewPassword(this.user).then(() => {
      this.dialog.closeAll();
      this.dialog.open(MessagePopupComponent, {data: {title: "Password updated!", text: "Log in with your new password!"}});
      this.isConfirm = false;
      this.user = {} as IUser;
    }).catch(() => {
      this.dialog.closeAll();
      this.dialog.open(MessagePopupComponent, {data: {title: "Error", text: "Check code and/or password"}});
    });
  }

  goToSignIn() {
    this.router.navigate(['/signIn']);
  }

}
