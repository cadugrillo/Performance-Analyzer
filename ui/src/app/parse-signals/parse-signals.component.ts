import { Component, OnInit, ViewChild } from '@angular/core';
import { SignalsService } from '../signals.service';

@Component({
  selector: 'app-parse-signals',
  templateUrl: './parse-signals.component.html',
  styleUrls: ['./parse-signals.component.css']
})
export class ParseSignalsComponent implements OnInit {

  @ViewChild('file') file: any

  constructor(private SignalsService: SignalsService) { }

  ngOnInit(): void {
  }

  importConfig() {
    this.file.nativeElement.click();
   }
  
   onFilesAdded() {
    const xlsx = this.file.nativeElement.files[0];
    this.file.nativeElement.value = "";
    const formData = new FormData();
    formData.append("Signals.xlsx", xlsx);

    this.SignalsService.parseSignals(xlsx).subscribe((data) => {
      console.log(data);
    });
   }

}
