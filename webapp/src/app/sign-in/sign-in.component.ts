import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MessagePopupComponent } from '../message-popup/message-popup.component';
import { WaitPopupComponent } from '../wait-popup/wait-popup.component';

import { IUser, CognitoService } from '../cognito.service';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css'],
})
export class SignInComponent {

  loading: boolean;
  user: IUser;

  constructor(private router: Router, public dialog: MatDialog,
              private cognitoService: CognitoService) {
    this.loading = false;
    this.user = {} as IUser;
  }

  public signIn(): void {
    if (this.user.email != null && this.user.password != null) {
      
      this.loading = true;
      this.dialog.open(WaitPopupComponent, {});
      this.cognitoService.signIn(this.user)
      .then(() => {
        this.router.navigate(['/home']).then(() => {window.location.reload();});
      }).catch(() => {
        this.loading = false;
        this.dialog.closeAll();
      });

    } else this.dialog.open(MessagePopupComponent, {data: {title: "Error", text: "Empty email or password!"}});
  }

  goToSignUp() {
    this.router.navigate(['/signUp']);
  }
  
}
