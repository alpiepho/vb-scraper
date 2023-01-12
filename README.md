

## vs Original

Added main.go to replace scraper1 and python spiders.  It is a single file/program, and should be easy to follow even without knowing the Golang syntax.

The output json file(s) are not the same as the previous instrumentinfo.json.  That file has a format (that I don't fully understand yet), that attempts to compress the info to fewer fields and characters, then the IOLS code uses that to reconstruct the url for images (and product data?).

Unfortunately, we need to verify that image urls in that format are still available.  There are several possible formats from code inspection and the limited set I tried were not available on k.com.

I am also thinking that it might be better to be in better sync with the latest k.com data.

## TODO
- [x] verify run times...50min
- [x] full outout, partial with entry for path and only png path
- [ ] option for headless
- [ ] option for debug

- [x] refactor with functions
- [x] better debug message control
- [ ] compare image data from catalog with instrumentinfo.json

- [ ] can we convert to old json format? NO
- [ ] how does IOLS reconstruct format? NO



## Golang version

The scraper1 utility was rewritten in Go and chromedp.  Go is a little
more difficult than python, but it seems to run a little faster.  Also,
chromedp will allow automating web flow and can run headless.

Install golang [here](https://golang.org/doc/install)

[chromedp](https://github.com/chromedp/chromedp) library.

`go get -u github.com/chromedp/chromedp`

`go run main.go`

or to gather all the output:

`go run main.go > results.txt 2>&1`

## Reference

https://confluence.it.keysight.com/pages/viewpage.action?pageId=85870460

IOLS code that uses instrumentinfo.json:

private List<string> IdentifyTemplateParameters(string template)
  


<!-- 
NOTE for getting product page link:

https://www.keysight.com/us/en/products.html
https://www.keysight.com/us/en/products/oscilloscopes.html
https://www.keysight.com/us/en/catalog/key-34568/oscilloscopes.html
https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/migrated/scene7/products/1a/PROD-2396122-01.png?wid=358&hei=202&fmt=png-alpha&resMode=sharp2&op_sharpen=1
data-parentvalue="DSOS054A High-Definition Oscilloscope: 500 MHz, 4 Analog Channels"

will redirect to product page":
https://www.keysight.com/us/en/product/DSOS054A 
-->

<!-- 
NOTE for getting original instrumentinfo.json
// The update scheme for Instrument Information is to have a jump station URL that returns a short JSON object with
// a timestamp for the latest Instrument Information and a URL where it can be found.
//
// The current jumpstation URL is: http://www.keysight.com/find/instrumentinfodata
//
// An example of the jumpstation JSON:
//
// {
//   "SchemaVersion":"1",
//   "LastUpdated":"2014/03/01 18:03",
//   "DataUrl":"http://www.keysight.com/main/redirector.jspx?action=ref&cname=EDITORIAL&ckey=2438582&cc=US&lc=eng"
// }
-->


<!-- 
NOTES: from meeting with IT

https://1kqa05t4gh.execute-api.us-west-2.amazonaws.com/kapi/krs/products/N9020B

...
"PRODUCT_HIGHLIGHTS": "...\r\n",
"PRODUCT_HIGHLIGHTS_VIEW": "<h3>...\n",
"TITLE": "N9020B MXA Signal Analyzer, 10 Hz to 50 GHx",
"ATTRIBUTES" : [

],
"IMAGE_PREFIX": "https://stagekeysight-h.assetsadobe.com/is/image",
"PRODUCT_IMAGE": "/content... .png",



https://www.keysight.com/us/en/product/N9020B/n9020b-mxa-signal-analyzer-multi-touch-10-hz-50-ghz.html
https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/x-series/n9020b/N9020B_26FL_1600x900.png

<div id="dynamicmedia_128090962" data-current-page="/content/keysight/us/en/products/spectrum-analyzers-signal-analyzers/x-series-signal-analyzers/n9020b-mxa-signal-analyzer-multi-touch-10-hz-50-ghz" data-page-locale="en" data-asset-path="/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/x-series/n9020b/N9020B_26FL-1600x900-set" data-asset-name="N9020B_26FL-1600x900-set" data-asset-type="imageset" data-viewer-type="ZoomViewer" data-viewer-path="/etc/dam/viewers/s7viewers/" data-imageserver="https://keysight-h.assetsadobe.com/is/image/" data-videoserver="https://gateway-na.assetsadobe.com/DMGateway/public/stagekeysight" data-contenturl="/" data-config="/conf/global/settings/dam/dm/presets/viewer/ImageSet_light|IMAGE_SET|false" data-wcmdisabled="" data-enablehd="never" data-linktarget="_self" class="s7dm-dynamic-media s7responsiveViewer s7zoomviewer s7mouseinput s7size_small s7device_landscape" style="height: auto;"><div id="dynamicmedia_128090962_container" class="s7container" data-description="Scene7ComponentHolder" data-component="Container" data-namespace="s7viewers" lang="en" mode="normal" role="region" aria-label="Zoom viewer" style="position: relative; height: 421px;"><div mode="normal" class="s7container s7innercontainer" id="dynamicmedia_128090962_container_inner"><div id="dynamicmedia_128090962_zoomView" class="s7zoomview" data-description="Scene7ComponentHolder" data-component="ZoomView" data-namespace="s7viewers" lang="en" tabindex="0" role="application" aria-roledescription="zoomable image" aria-describedby="dynamicmedia_128090962_zoomView_hint" style="width: 578px; height: 325px; overflow: hidden;" cursortype="reset"><span id="dynamicmedia_128090962_zoomView_hint" style="display: none;">Use + and - keys to zoom in and out, escape key to reset, arrow keys to change image in reset state or move the zoomed portion of the image</span><div style="z-index: 0; position: absolute; left: -578px; width: 1734px; height: 325px;"><div style="z-index: 900; position: absolute; width: 0px; height: 400px; left: 0px;"><canvas width="1" height="1" style="position: absolute; width: 0px; height: 400px;"></canvas></div><div style="z-index: 900; position: absolute; width: 578px; height: 325px; left: 0px;"><canvas width="867" height="487" style="position: absolute; width: 578px; height: 325px;"></canvas></div><div style="z-index: 910; position: absolute; width: 578px; height: 325px; left: 578px;"><canvas width="867" height="487" style="position: absolute; width: 578px; height: 325px;"></canvas></div><div style="z-index: 920; position: absolute; width: 578px; height: 325px; left: 1156px;"><canvas width="867" height="487" style="position: absolute; width: 578px; height: 325px;"></canvas></div></div></div><div id="dynamicmedia_128090962_divcontainer" style="position: absolute; width: 578px; top: 325px; height: 0px; z-index: 1;"><div id="dynamicmedia_128090962_zoomInButton" class="s7button s7zoominbutton" data-description="Scene7ComponentHolder" data-component="ZoomInButton" data-namespace="s7viewers" lang="en" role="button" aria-label="Zoom In" state="disabled" aria-disabled="true"></div><div id="dynamicmedia_128090962_zoomOutButton" class="s7button s7zoomoutbutton" data-description="Scene7ComponentHolder" data-component="ZoomOutButton" data-namespace="s7viewers" lang="en" role="button" aria-label="Zoom Out" state="up" tabindex="0"></div><div id="dynamicmedia_128090962_zoomResetButton" class="s7button s7zoomresetbutton" data-description="Scene7ComponentHolder" data-component="ZoomResetButton" data-namespace="s7viewers" lang="en" role="button" aria-label="Reset Zoom" state="up" tabindex="0"></div><div id="dynamicmedia_128090962_fullScreenButton" class="s7button s7fullscreenbutton" data-description="Scene7ComponentHolder" data-component="FullScreenButton" data-namespace="s7viewers" lang="en" tabindex="0" role="button" aria-label="Full Screen" state="up" selected="false"></div></div><div id="dynamicmedia_128090962_swatches" class="s7swatches" data-description="Scene7ComponentHolder" data-component="Swatches" data-namespace="s7viewers" lang="en" tabindex="-1" role="listbox" aria-label="swatches" style="width: 578px; height: 96px; position: absolute;"><div style="width: 340px; height: 66px; position: absolute; overflow: hidden; left: 119px; top: 15px;"><div id="dynamicmedia_128090962_swatches_listbox" style="position: absolute; width: 330px; height: 66px; left: 0px; top: 0px; transform: translateZ(0px);"><div data-namespace="s7viewers" class="s7thumbcell" tabindex="-1" role="option" aria-selected="false" swindex="0" aria-setsize="5" aria-posinset="1" style="margin: 0px; position: absolute; left: 5px; top: 5px;"><div data-namespace="s7viewers" class="s7thumb" state="default" style="width: 56px; height: 56px; background-image: url(&quot;https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/x-series/n9020b/N9020B_26FL_1600x900.png?fit=constrain,1&amp;wid=56&amp;hei=56&amp;fmt=jpg&quot;);"><div data-namespace="s7viewers" class="s7thumboverlay" type="image"></div></div></div><div data-namespace="s7viewers" class="s7thumbcell" tabindex="0" role="option" aria-selected="true" swindex="1" aria-setsize="5" aria-posinset="2" style="margin: 0px; position: absolute; left: 71px; top: 5px;"><div data-namespace="s7viewers" class="s7thumb" state="selected" style="width: 56px; height: 56px; background-image: url(&quot;https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/x-series/n9020b/N9020B_26.5GHz_6_1600x900.png?fit=constrain,1&amp;wid=56&amp;hei=56&amp;fmt=jpg&quot;);"><div data-namespace="s7viewers" class="s7thumboverlay" type="image"></div></div></div><div data-namespace="s7viewers" class="s7thumbcell" tabindex="-1" role="option" aria-selected="false" swindex="2" aria-setsize="5" aria-posinset="3" style="margin: 0px; position: absolute; left: 137px; top: 5px;"><div data-namespace="s7viewers" class="s7thumb" state="default" style="width: 56px; height: 56px; background-image: url(&quot;https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/x-series/n9020b/N9020B_26_FL_1600x900.png?fit=constrain,1&amp;wid=56&amp;hei=56&amp;fmt=jpg&quot;);"><div data-namespace="s7viewers" class="s7thumboverlay" type="image"></div></div></div><div data-namespace="s7viewers" class="s7thumbcell" tabindex="-1" role="option" aria-selected="false" swindex="3" aria-setsize="5" aria-posinset="4" style="margin: 0px; position: absolute; left: 203px; top: 5px;"><div data-namespace="s7viewers" class="s7thumb" state="default" style="width: 56px; height: 56px; background-image: url(&quot;https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/x-series/n9020b/N9020B_26_1600x900.png?fit=constrain,1&amp;wid=56&amp;hei=56&amp;fmt=jpg&quot;);"><div data-namespace="s7viewers" class="s7thumboverlay" type="image"></div></div></div><div data-namespace="s7viewers" class="s7thumbcell" tabindex="-1" role="option" aria-selected="false" swindex="4" aria-setsize="5" aria-posinset="5" style="margin: 0px; position: absolute; left: 269px; top: 5px;"><div data-namespace="s7viewers" class="s7thumb" state="default" style="width: 56px; height: 56px; background-image: url(&quot;https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/prd/scopes-analyzers-meters/spectrumanalyzers/N9020B_26_1600x900.png?fit=constrain,1&amp;wid=56&amp;hei=56&amp;fmt=jpg&quot;);"><div data-namespace="s7viewers" class="s7thumboverlay" type="image"></div></div></div></div></div></div></div></div></div>
-->

<!--
JIRA from IT: https://jira.it.keysight.com/browse/KESWP-73

https://dragon.is.keysight.com/aemlinks/reports/loader.shtml

4 sheets that get data from AMAZON

AEM Reports
The following Excel files are connected to reports that run regularly. Once the spreadsheet is downloaded, you can use the Refresh All icon on the Data ribbon to fetch the latest report data.

Download Excel: AEM Pages Report (aemmap)
Highlighted cells indicate where translations are out of date compared to the English. Does not include dynamic product pages or pages in Resources/Events/Software Details.
spreadsheet template last updated: Wednesday, 16-Feb-2022 23:02:46 UTC
Download Excel: AEM Asset Data Report
Includes all indexed pubkey assets for pre-sales and post-sales content types.
spreadsheet template last updated: Wednesday, 16-Feb-2022 00:03:29 UTC
Download Excel: Model Relations Report
Includes all published model numbers in the web hierarchy, with information about relations configured in PIM for accessories, software, upgrade products, quote cross sell products, and replacement products.
spreadsheet template last updated: Monday, 21-Mar-2022 13:43:16 UTC
Download Excel: Model Fragments Report
Includes all authored model numbers in the web hierarchy, with information about experience fragments and product images.
spreadsheet template last updated: Thursday, 24-Mar-2022 20:03:53 UTC

Model Relations Report has 10000 Models, descriptions, link to info, but no image
Model Fragment Report has 3700+ Models, descriptions, link to info, and images...b3ut 3702 is less than 4021 found from web or 3991 from old instrumentinfo.json
-->