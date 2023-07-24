package mailer

import "github.com/supabase/gotrue/internal/conf"

func BaseLayout(config *conf.GlobalConfiguration) string {
	return `<!DOCTYPE html>
<html lang="en" xmlns:v="urn:schemas-microsoft-com:vml">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <style>
    td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
  </style>
  <![endif]-->
  <style>
    @media (max-width: 600px) {
      .sm-ml-5px {
        margin-left: 5px !important
      }
      .sm-mr-5px {
        margin-right: 5px !important
      }
      .sm-w-1-3 {
        width: 33.333333% !important
      }
      .sm-w-2-3 {
        width: 66.666667% !important
      }
      .sm-w-280px {
        width: 280px !important
      }
      .sm-w-320px {
        width: 320px !important
      }
      .sm-w-80px {
        width: 80px !important
      }
      .sm-px-4 {
        padding-left: 16px !important;
        padding-right: 16px !important
      }
      .sm-pl-5px {
        padding-left: 5px !important
      }
    }
  </style>
</head>
<body style="margin: 0; width: 100%; padding: 0; -webkit-font-smoothing: antialiased; word-break: break-word">
  <div role="article" aria-roledescription="email" aria-label lang="en">
    <div style="background-color: #fff">
      <table align="center" cellpadding="0" cellspacing="0" role="presentation">
        <tr>
          <td style="width: 800px; max-width: 100%">
            <div style="margin-top: 16px; text-align: center">
              <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/AladinMall.png" alt="aladinmall logo" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0">
            </div>
          </td>
        </tr>
      </table>
      <table align="center" cellpadding="0" cellspacing="0" role="presentation">
        <tr>
          <td style="width: 800px; max-width: 100%;">
            <div style="display: flex; align-items: center; justify-content: center">
              <div class="sm-w-320px" style="margin-top: 46px; height: auto; width: 800px; border-radius: 10px; border: 1px solid #d9d9d9; background-color: #fff; color: #000">
                <div class="sm-px-4" style="padding-left: 90px; padding-right: 90px; padding-top: 46px">
                  <p style="margin: 0; padding-bottom: 8px; font-size: 16px; font-weight: 400; line-height: 22px">Hai <b>Aladiners,</b></p>
                  <p style="margin: 0; padding-bottom: 8px; font-size: 14px; line-height: 22px"></p>

                  {{content}}

                  <p style="margin-bottom: 8px; margin-top: 30px; padding: 0; font-size: 14px; line-height: 22px"></p>
                  <p style="margin-top: 8px; font-size: 14px; line-height: 22px">Salam hangat,</p>
                  <p style="margin-bottom: 8px; margin-top: 30px; padding: 0; font-size: 14px; line-height: 22px;"><b>Tim AladinMall</b></p>
                  <table align="center" style="margin-left: auto; margin-right: auto" cellpadding="0" cellspacing="0" role="presentation">
                    <tr>
                      <td style="max-width: 100%;">
                        <div style="margin-bottom: 8px; display: flex;">
                          <a href="` + config.SiteURL + `" target="_blank" style="cursor: pointer">
                            <button style="height: 30px; width: 116px; cursor: pointer; border-radius: 8px; border: 1px solid #e64325; background-color: #E64325; color: #fff">Ke AladinMall</button>
                          </a>
                        </div>
                      </td>
                    </tr>
                  </table>
                  <p style="margin: 0 0 30px; padding: 0; text-align: center; font-size: 10px; font-style: italic">Klik "Unsubscribe" untuk berhenti menerima email seperti ini dari
                    kami lagi. Kami akan sangat merindukan kehadiran Anda, tetapi kami menghormati keputusan Anda.
                    Terima kasih.</p>
                  <div style="margin-bottom: 12px; width: fit-content">
                    <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Rectangle-Body.png" alt="horizontal line" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
                  </div>
                </div>
                <table align="center" style="margin-bottom: 40px; margin-left: auto; margin-right: auto" cellpadding="0" cellspacing="0" role="presentation">
                  <tr>
                    <td class="sm-w-1-3">
                      <div class="sm-ml-5px sm-mr-5px" style="margin-right: 20px; padding: 0">
                        <p style="margin: 0; padding: 0; font-size: 12px; line-height: 20px; color: #E64325"><b>Part of</b></p>
                        <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/aladinmall-logo.png" alt="MNC logo" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
                      </div>
                    </td>
                    <td class="sm-w-2-3" style="width: 518px">
                      <div style="display: flex; flex-direction: column">
                        <p style="font-size: 10px; font-weight: 400; line-height: 12px">AladinMall merupakan toko belanja online terlengkap dan
                          terpercaya.
                          Menyediakan beragam pilihan produk kebutuhan sehari-hari dengan jaminan harga termurah
                          dan kualitas terbaik. Layanan pengiriman luas serta. kemudahan pembayaran bagi seluruh
                          pelanggan.</p>
                      </div>
                    </td>
                  </tr>
                </table>
              </div>
            </div>
          </td>
        </tr>
      </table>
      <table align="center" cellpadding="0" cellspacing="0" role="presentation">
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <div style="margin-bottom: 25px; margin-top: 33px; text-align: center">
              <p style="margin: 0; padding: 0; font-size: 12px; font-weight: 400; font-style: italic; line-height: 20px; color: #8C8C8C">*Email dikirim secara otomatis.
                Harap jangan mengirim balasan ke email ini.
              </p>
            </div>
          </td>
        </tr>
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <table align="center" style="margin-bottom: 25px;" cellpadding="0" cellspacing="0" role="presentation">
              <tr>
                <td style="width: 22px; max-width: 100%; padding-right: 27px; text-align: center">
                  <a href="https://www.tiktok.com/@aladinmall" target="_blank">
                    <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Tiktok.png" alt="tiktok logo" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
                  </a>
                </td>
                <td style="width: 22px; max-width: 100%; padding-right: 27px; text-align: center;">
                  <a href="https://www.facebook.com/aladinmall.id" target="_blank">
                    <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Facebook.png" alt="facebook logo" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
                  </a>
                </td>
                <td style="width: 22px; max-width: 100%; padding-right: 27px; text-align: center;">
                  <a href="https://www.youtube.com/@aladinmall" target="_blank">
                    <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Youtube.png" alt="youtube logo" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
                  </a>
                </td>
                <td style="width: 22px; max-width: 100%; text-align: center;">
                  <a href="https://www.instagram.com/aladinmall.id/" target="_blank">
                    <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Instagram.png" alt="instagram logo" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
                  </a>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <div style="text-align: center;">
              <p style="margin: 0; padding: 0; font-size: 12px; font-weight: 400; line-height: 20px; color: #8C8C8C;">MNC Center lt. 20 Jl. Kebon Sirih
                No.17-19, Jakarta Pusat 10340</p>
            </div>
          </td>
        </tr>
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <div style="margin-bottom: 30px; margin-top: 30px; width: 100%;">
              <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Rectangle-Footer.png" alt="horizontal line" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
            </div>
          </td>
        </tr>
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <div style="text-align: center;">
              <p style="margin: 0; padding-top: 8px; font-size: 12px; font-weight: 400; line-height: 20px; color: #8C8C8C">Senin - Sabtu : (021) 390 6001 |
                Whatsapp: +62 811 113 8080</p>
              <p style="margin: 0; padding-top: 8px; font-size: 12px; color: #8C8C8C;">Minggu : (hanya email)
                cs-aladinmall@misteraladin.com</p>
            </div>
          </td>
        </tr>
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <div style="margin-bottom: 30px; margin-top: 30px; width: 100%;">
              <img src="https://user-api-stg.aladinmall.id/storage/v1/object/public/web/public/email/Rectangle-Footer.png" alt="horizontal line" style="max-width: 100%; vertical-align: middle; line-height: 1; border: 0;">
            </div>
          </td>
        </tr>
        <tr>
          <td class="sm-w-320px" style="width: 800px; max-width: 100%">
            <div style="text-align: center;">
              <p style="margin-bottom: 28px; padding: 0; font-size: 12px; font-weight: 400; line-height: 20px; color: #8C8C8C">&copy; 2024 AladinMall. PT MNC
                ALADIN INDONESIA All rights reserved.</p>
            </div>
          </td>
        </tr>
      </table>
    </div>
  </div>
</body>
</html>`
}
