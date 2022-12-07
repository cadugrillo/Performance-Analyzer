import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import {animate, state, style, transition, trigger} from '@angular/animations';
import { Subscription } from 'rxjs';
import { IMqttMessage } from "ngx-mqtt";
import { MqttClientService } from '../../services/mqttClient.service';
import { SignalsService, AnalyzedData } from '../../services/signals.service';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
import { MessagePopupComponent } from '../message-popup/message-popup.component';
import { WaitPopupComponent } from '../wait-popup/wait-popup.component';
import { MatDialog } from '@angular/material/dialog';
import { saveAs } from "file-saver";
import { env } from 'process';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-mqtt-client',
  templateUrl: './mqtt-client.component.html',
  styleUrls: ['./mqtt-client.component.css'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({height: '0px', minHeight: '0'})),
      state('expanded', style({height: '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class MqttClientComponent implements OnInit, OnDestroy {

  @ViewChild('json') json: any;
  messages: IMqttMessage[] = [];
  recTS: string[] = [];
  subscription!: Subscription;
  topic: string = '';
  maxCapturedMessages = 10000;
  maxMessagesPerExport = 1000000;
  running: boolean = false;
  dataSource!: MatTableDataSource<IMqttMessage>;
  telegramsToAnalize!: Object;
  telegramsImportStatus: string = ""
  analyzedData!: AnalyzedData;
  analysisRunning: boolean = false;
  TsInterval!: number;
  checkbox: boolean = true;
  broker_url: string = 'tcp://' + environment.mqtt.server + ":" + environment.mqtt.pub_port;
  qos: 0 | 1 | 2 = 2;
  

  columnsToDisplay = ['Topic', 'Timestamp'];
  columnsToDisplayWithExpand = [...this.columnsToDisplay, 'expand'];
  expandedElement!: IMqttMessage | null;

  @ViewChild(MatPaginator) paginator!: MatPaginator;

  constructor(private mqttClientService: MqttClientService, public dialog: MatDialog,
              private SignalsService: SignalsService) {}

  ngOnInit(): void {
    this.dataSource = new MatTableDataSource();
  }

  ngOnDestroy(): void {
    
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
    
    this.running = false;
    this.messages = [];
    this.recTS = [];  
    this.analysisRunning = false;
}

  subscribeToTopic(topic: string) {
    if (this.topic != '' && this.maxCapturedMessages >= 1 && this.maxCapturedMessages <= 1000000) {
      this.running = true;
      this.subscription = this.mqttClientService.topic(topic, this.qos).subscribe((data: IMqttMessage) => {
        //console.log('Initial time:'+this.getTimestamp("display"));
        this.messages.push(data);
        this.recTS.push(this.getTimestamp("display"));
        this.dataSource = new MatTableDataSource(this.messages);
        this.dataSource.paginator = this.paginator;
        if (this.messages.length >= this.maxCapturedMessages) {
          this.unsubscribeTopic();
        }
        //console.log('Final time:'+this.getTimestamp("display"));
      });
    } else this.dialog.open(MessagePopupComponent, {data: {title: "Error", text: "Topic field cannot be empty and/or number of captured messages should be between 1 and 1000000"}});
  }

  unsubscribeTopic() {
    this.subscription.unsubscribe();
    this.running = false; 
  }

  clearData() {
    this.dataSource = new MatTableDataSource();
    this.messages = [];
    this.recTS = [];
  }

  toString(payload: Object): string {
    return JSON.stringify(JSON.parse(payload.toString()), null, 2);
  }

  toStringData(parsedSignals: Object): string {
    return JSON.stringify(parsedSignals, null, 4);
  }

  getTimestamp(format: string): string {

    switch(format) {
      case "display":
        var today = new Date();
        var date = today.getFullYear()+'-'+(today.getMonth()+1)+'-'+today.getDate();
        var time = today.getHours() + ":" + today.getMinutes() + ":" + today.getSeconds()+ ":" + today.getMilliseconds();
        var dateTime = date+' '+time;
        return dateTime;
      case "file":
        var today = new Date();
        var date = today.getFullYear()+'_'+(today.getMonth()+1)+'_'+today.getDate();
        var time = today.getHours() + "_" + today.getMinutes() + "_" + today.getSeconds()+ "_" + today.getMilliseconds();
        var dateTime = date+'_'+time;
        return dateTime;
    }
    return "err"
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  exportMessages() {

    var j = 0;
    var k = 0;
    let exportedData: ExportedData[] = [];

    for (var i = 0; i < this.messages.length; i++ ){

      let exportedDataItem = new ExportedData;
      exportedDataItem.topic = this.messages[i].topic;
      exportedDataItem.payload = JSON.parse(this.messages[i].payload.toString());
      exportedData.push(exportedDataItem);
      j = j + 1;

      if (j == this.maxMessagesPerExport || i == this.messages.length - 1) {
        saveAs(new Blob([JSON.stringify(exportedData, null, 2)], { type: 'JSON' }), 'messages_'+this.getTimestamp("file")+'_part_'+k+'.json');
        j = 0;
        k = k +1;
        exportedData = [];
      }
    }
    return 
  }

  wrapMessages(): ExportedData[] {

    let exportedData: ExportedData[] = [];
    this.messages.forEach((message) => {
      let exportedDataItem = new ExportedData
      exportedDataItem.topic = message.topic;
      exportedDataItem.payload = JSON.parse(message.payload.toString());
      exportedData.push(exportedDataItem);
    });
    return exportedData
  }

  analyzeFromCapture(tsInterval: number) {
    this.dialog.open(WaitPopupComponent, {});
    this.analysisRunning = true;

    if (this.checkbox) {
      this.SignalsService.analizeTelegramsData(this.wrapMessages(), tsInterval).subscribe((data) => {
        this.analyzedData = data as AnalyzedData;
        this.dialog.closeAll();
        this.analysisRunning = false;
        this.dialog.open(MessagePopupComponent, {data: {title: "Analysis Finished", text: "Check Results!"}});
      });
    } else {
      this.SignalsService.analizeTelegramsDbusData(this.wrapMessages(), tsInterval).subscribe((data) => {
        this.analyzedData = data as AnalyzedData;
        this.dialog.closeAll();
        this.analysisRunning = false;
        this.dialog.open(MessagePopupComponent, {data: {title: "Analysis Finished", text: "Check Results!"}});
      });
    }

  }

  importTelegramsToAnalyze() {
    this.json.nativeElement.click();
   }

  onTelegramasAdded(tsInterval: number) {
    this.dialog.open(WaitPopupComponent, {});
    this.analysisRunning = true;
    const jsonfile = this.json.nativeElement.files[0];
    this.json.nativeElement.value = "";

    if (this.checkbox) {
      this.SignalsService.analizeTelegramsData(jsonfile, tsInterval).subscribe((data) => {
        this.analyzedData = data as AnalyzedData;
        this.dialog.closeAll();
        this.analysisRunning = false;
        this.dialog.open(MessagePopupComponent, {data: {title: "Analysis Finished", text: "Check Results!"}});
      });
    } else {
      this.SignalsService.analizeTelegramsDbusData(jsonfile, tsInterval).subscribe((data) => {
        this.analyzedData = data as AnalyzedData;
        this.dialog.closeAll();
        this.analysisRunning = false;
        this.dialog.open(MessagePopupComponent, {data: {title: "Analysis Finished", text: "Check Results!"}});
      });
    }
  }

  exportAnalyzedData() {
    return saveAs(new Blob([JSON.stringify(this.analyzedData, null, 2)], { type: 'JSON' }), 'AnalyzedData'+this.getTimestamp("file")+'.json');
  }

  clearAnalyzedData() {
    this.telegramsToAnalize = {};
    this.analyzedData = new AnalyzedData();
  }

}

class ExportedData {
  topic!: Object
  payload!: Object
}