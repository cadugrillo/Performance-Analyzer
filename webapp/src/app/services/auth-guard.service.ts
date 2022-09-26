import { Injectable } from '@angular/core';
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { CognitoService } from './cognito.service';

@Injectable({
  providedIn: 'root'
})

export class AuthGuardService implements CanActivate {

  constructor(private cognitoService: CognitoService, 
              private router: Router) {}


  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Promise<boolean> {
   
    return new Promise((resolve, reject) => {
      
      return resolve (this.cognitoService.isAuthenticated())

    });
  }
}
