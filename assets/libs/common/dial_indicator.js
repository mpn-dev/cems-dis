class DialIndicator {
  constructor(canvas, options = {}) {
    this.canvas = canvas;
    this.ctx = canvas.getContext("2d");
    this.angleOffset = 360 - 90;
    this.centerOffsetY = options.centerOffsetY ? options.centerOffsetY : 0;
    this.centerOffsetX = options.centerOffsetX ? options.centerOffsetX : 0;
    this.min = options.min ? options.min : 0;
    this.max = options.max ? options.max : 100;
    this.major = options.major ? options.major : null;
    this.minor = options.minor ? options.minor : null;
    this.radius = options.dialRadius ? options.dialRadius : Math.min(this.canvas.width / 2, this.canvas.height / 2) - 10;
    this.startAngle = options.dialStartAngle ? options.dialStartAngle : 0;
    this.endAngle = options.dialEndAngle ? options.dialEndAngle : 360;
    this.dialWidth = options.dialWidth ? options.dialWidth : 15;
    this.dialColor = options.dialColor ? options.dialColor : "#000000";
    this.barWidth = options.barWidth ? options.barWidth : 13;
    this.barColor = options.barColor ? options.barColor : "#00ff00";
    this.labelFont = options.labelFont ? options.labelFont : "13px serif";
    this.labelColor = options.labelColor ? options.labelColor : "#ff0000";
    this.labelOffsetY = options.labelOffsetY ? options.labelOffsetY : 0;
    this.labelValue = options.labelValue ? options.labelValue : null;
    this.labelOffsetX = options.labelOffsetX ? options.labelOffsetX : 0;
    this.gridFont = options.gridFont ? options.gridFont : "12px serif";
    this.dials = options.dials ? options.dials : null;
    this.needleColor = options.needleColor ? options.needleColor : "rgba(255,0,0,0.8)";
    this.x = (canvas.width / 2) + this.centerOffsetX;
    this.y = (canvas.height / 2) + this.centerOffsetY;
    this.ctx.globalAlpha = 1.0;

    this.getBarColor = function(value) {
      if (typeof this.barColor === "function") {
        return this.barColor(value);
      } else {
        return this.barColor;
      }
    }

    this.getLabelValue = function(value) {
      if(typeof this.labelValue === "function") {
        return this.labelValue(value);
      } else if(value == null) {
        return value;
      } else {
        if(Math.round(value, 0) == value) {
          return Math.round(value, 0);
        } else {
          return value.toFixed(1);
        }
      }
    }

    this.gridLabelAt = function(value) {
      return value;
    }
  }

  _pointAt(angle, radius) {
    return {x: this.x + radius * Math.sin((angle / 360) * 2 * Math.PI), y: this.y - radius * Math.cos((angle / 360) * 2 * Math.PI)};
  }

  _generateGrids(minValue, maxValue) {
    var range = maxValue - minValue;
    var exp = range < 1 ? Math.floor(Math.log10(range * 100000)) - 5 : Math.floor(Math.log10(range));
    var mag = 10 ** exp;

    var rpm = range / (mag * 1.0);
    var rpm10 = (rpm >= 8.5 ? 10 : rpm) * 10;
    var category = 0, major = 0, minor = 0;

    if((range >= 0.8) && (range <= 6)) {
      if((range >= 1.5) && (range <= 6)) {
        category = 1;
        major = 1;
        minor = major / 5;
      } else if((range >= 0.8) && (range <= 1.5)) {
        category = 2;
        major = 0.2;
        minor = major / 4;
      }
    } else if(rpm10 <= 11) {
      category = 3;
      major = mag * 0.1;
      minor = major / 5;
    } else if((rpm10 <= 18) && (range <= 24)) {
      category = 4;
      major = mag * 0.2;
      minor = major / 4;
    } else if(((rpm10 % 3) == 0) && ((rpm10 % 9) == 0)) {
      category = 5;
      major = mag * 0.3;
      minor = major / 3;
    } else if((rpm10 % 5) == 0) {
      if((rpm10 % 10) != 0) {
        category = 6;
        major = mag * 0.5;
        minor = major / 5;
      } else {
        category = 7;
        major = mag * 1.0;
        minor = major / 5;
      }
    } else if(rpm10 == 50) {
      category = 8;
      major = mag * 1.0;
      minor = major / 5;
    } else if(rpm10 == 10) {
      category = 9;
      major = mag * 0.1;
      minor = major / 5;
    } else if((rpm10 % 2) == 0) {
      category = 10;
      major = mag * 0.2;
      minor = major / 4;
    } else {
      category = 11;
      major = mag * 0.1;
      minor = major / 5;
    }

    return {major: major, minor: minor};
  }

  _clearCanvas() {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
  }

  _valueToAngle(value) {
    var angleRange = this.endAngle + (this.endAngle < this.startAngle ? 360 : 0) - this.startAngle;
    var valueRange = this.max - this.min;
    return this.startAngle + (((value - this.min) / valueRange) * angleRange);
  }

  _drawCircle(x, y, radius, fillStyle) {
    this.ctx.fillStyle = fillStyle;
    this.ctx.beginPath();
    this.ctx.arc(x, y, radius, 0, 2 * Math.PI, false);
    this.ctx.fill();
  }

  _drawArc(startAngle, endAngle, radius, width, color) {
    this.ctx.lineWidth = width;
    this.ctx.strokeStyle = color;
    this.ctx.beginPath();
    this.ctx.arc(this.x, this.y, radius, ((this.angleOffset + startAngle) / 360) * 2 * Math.PI, ((this.angleOffset + endAngle) / 360) * 2 * Math.PI, false);
    this.ctx.stroke();
  }

  _drawGrid(angle, radius, length, color) {
    var p1 = this._pointAt(angle, radius);
    var p2 = this._pointAt(angle, radius - length);
    this.ctx.lineWidth = 1;
    this.ctx.strokeStyle = color;
    this.ctx.beginPath();
    this.ctx.moveTo(p1.x, p1.y);
    this.ctx.lineTo(p2.x, p2.y);
    this.ctx.stroke();
  }

  _drawBackground() {
    this._drawCircle(this.x, this.y, this.radius, "#ffffff");
    this._drawArc(0, 360, this.radius * 0.97, this.radius * 0.06, "#e0e0e0");
    this._drawArc(0, 360, this.radius, 1, "#000000");
  }

  _drawDials() {
    if(this.dials != null) {
      var width = this.radius * 0.14;
      var radius = (this.radius * 0.9) - (width / 2);

      this.dials.forEach(function(d) {
        this._drawArc(this._valueToAngle(d.min), this._valueToAngle(d.max), radius, width, d.color);
      }, this)
    }
  }

  _drawGrids() {
    var markGrid = function(value, radius, length, showLabel) {
      var angle = this._valueToAngle(value);
      var p1 = this._pointAt(angle, radius);
      var p2 = this._pointAt(angle, radius - length);
      this.ctx.moveTo(p1.x, p1.y);
      this.ctx.lineTo(p2.x, p2.y);

      if(showLabel) {
        if(!((this.startAngle == 0) && (this.endAngle == 360) && (value == this.max))) {
          var p3 = this._pointAt(angle, radius - length - 10);
          var text = "" + this.gridLabelAt(value);
          var size = this.ctx.measureText(text);
          this.ctx.fillText(text, p3.x - (size.width / 2), p3.y + 4);
        }
      }
    }

    var applyGrids = function(step, radius, size, width, showLabel) {
      if(showLabel) {
        this.ctx.font = this.gridFont;
        this.ctx.fillStyle = "#000000";
      }

      this.ctx.lineWidth = width;
      this.ctx.strokeStyle = "#000000";
      this.ctx.beginPath();

      var baseValue = (this.min / step) == Math.floor((this.min / step)) ? this.min : (step * Math.ceil(this.min / step));
      if(this.min < baseValue) {
        markGrid.call(this, this.min, radius, size, showLabel);
      }

      var val = baseValue;
      var cnt = 0;
      while(val <= this.max) {
        markGrid.call(this, val, radius, size, showLabel);
        cnt = cnt + 1;
        val = Math.round((baseValue + cnt * step) * 1000000) / 1000000;
      }

      this.ctx.stroke();
    }

    var defaultGrids = this._generateGrids(this.min, this.max);
    var grids = {
      major: this.major != null ? this.major : defaultGrids.major, 
      minor: this.minor != null ? this.minor : defaultGrids.minor, 
    };
    var radius = this.radius * 0.9;
    applyGrids.call(this, grids.minor, radius, this.radius * 0.07, 0.5);
    applyGrids.call(this, grids.major, radius, this.radius * 0.14, 2, true);
  }

  _drawValue() {
    if(this.labelFont) {this.ctx.font = this.labelFont;}
    if(this.labelColor) {this.ctx.fillStyle = this.labelColor;}
    var temp = this.getLabelValue(this.value);
    var text = "" + (temp == null ? "" : temp);
    var size = this.ctx.measureText(text);
    this.ctx.fillText(text, this.x - (size.width / 2) + this.labelOffsetX, this.y + this.labelOffsetY);
  }

  _drawNeedle() {
    if(this.value == null) {
      return;
    }

    var value = this.value >= this.min ? (this.value <= this.max ? this.value : this.max) : this.min;

    this.ctx.strokeStyle = this.needleColor;
    this.ctx.lineWidth = 3;

    this.ctx.beginPath();
    this.ctx.moveTo(this.x, this.y);
    var point = this._pointAt(this._valueToAngle(value), this.radius * 0.67);
    this.ctx.lineTo(point.x, point.y);
    this.ctx.stroke();

    this.ctx.fillStyle = "#4295c4";
    this.ctx.strokeStyle = "#000000";
    this.ctx.lineWidth = 0.8;
    this.ctx.beginPath();
    this.ctx.arc(this.x, this.y, this.radius * 0.07, 0, 2 * Math.PI, false);
    this.ctx.fill();
    this.ctx.stroke();
  }

  render() {
    this._clearCanvas();
    this._drawBackground();
    this._drawDials();
    this._drawGrids();
    this._drawValue();
    this._drawNeedle();
  }

  setValue(val) {
    this.value = val;
    this.render();
  }
}


class Gauge extends DialIndicator {}


class Compass extends DialIndicator {
  constructor(canvas, options = {}) {
    super(canvas, options);
    this.min = 0;
    this.max = 360;
    this.dialStartAngle = 0;
    this.dialEndAngle = 360;

    this.gridLabelAt = function(value) {
      if((value % 90) == 0) {
        return ["N", "E", "S", "W"][(value /  90)];
      } else {
        return "";
      }
    }
  }
}
