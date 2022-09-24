import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MessagePopupComponent } from '../message-popup/message-popup.component';
import { WaitPopupComponent } from '../wait-popup/wait-popup.component';

import { IUser, CognitoService } from '../cognito.service';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.css'],
})
export class SignUpComponent {

  loading: boolean;
  isConfirm: boolean;
  user: IUser;

  constructor(private router: Router, public dialog: MatDialog,
              private cognitoService: CognitoService) {
    this.loading = false;
    this.isConfirm = false;
    this.user = {} as IUser;
  }

  public signUp(): void {
    if (this.user.email != null && this.user.password != null) {

      this.loading = true;
      this.dialog.open(WaitPopupComponent, {});
      this.cognitoService.signUp(this.user)
      .then(() => {
        this.loading = false;
        this.dialog.closeAll();
        this.isConfirm = true;
      }).catch(() => {
        this.loading = false;
        this.dialog.closeAll();
      });
      
    } else this.dialog.open(MessagePopupComponent, {data: {title: "Error", text: "Empty email or password!"}});
  }

  public confirmSignUp(): void {
    this.loading = true;
    this.dialog.open(WaitPopupComponent, {});
    this.cognitoService.confirmSignUp(this.user).then(() => {
      this.cognitoService.signIn(this.user).then(() => {
        this.router.navigate(['/profile']).then(() => {window.location.reload();});
      });    
    }).catch(() => { 
      this.loading = false;
      this.dialog.closeAll();
    }); 
  }

  goToSignIn() {
    this.router.navigate(['/signIn']);
  }

}
